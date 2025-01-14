//package simulator
//
//import (
//	"DES-go/schedulers/types"
//	"fmt"
//)
//
////type GPUType string
////type GPUID int
//
//type GPU struct {
//	gpuID          types.GPUID
//	gpuType        types.GPUType
//	gpuSourMilli   types.GPUSourMilli
//	gpuUsedMilli   types.GPUUsedMilli
//	gpuRemainMilli types.GPURemainMilli
//}
//
//func NewGPU(gpuID types.GPUID, gpuType types.GPUType, gpuSourMilli types.GPUSourMilli, gpuUsedMilli types.GPUUsedMilli, gpuRemainMilli types.GPURemainMilli) *GPU {
//	return &GPU{
//		gpuID:          gpuID,
//		gpuType:        gpuType,
//		gpuSourMilli:   gpuSourMilli,
//		gpuUsedMilli:   gpuUsedMilli,
//		gpuRemainMilli: gpuRemainMilli,
//	}
//}
//
//func (g *GPU) ID() types.GPUID {
//	return g.gpuID
//}
//
//func (g *GPU) Type() types.GPUType {
//	return g.gpuType
//}
//
//func (g *GPU) String() string {
//	return fmt.Sprintf("gpu=[ID=%d, Type=%s]", g.ID(), g.Type())
//}
//
//func (g *GPU) GpuSourMilli() types.GPUSourMilli {
//	return g.gpuSourMilli
//}
//
//func (g *GPU) GpuUsedMilli() types.GPUUsedMilli {
//	return g.gpuUsedMilli
//}
//
//func (g *GPU) GpuRemainMilli() types.GPURemainMilli {
//	return g.gpuRemainMilli
//}

// ------------------------------------------------------------------------------------
package simulator

import (
	"DES-go/schedulers/types"
	"fmt"
)

//type GPUType string
//type GPUID int

type GPU struct {
	gpuID   types.GPUID
	gpuType types.GPUType
}

func NewGPU(gpuID types.GPUID, gpuType types.GPUType) *GPU {
	return &GPU{
		gpuID:   gpuID,
		gpuType: gpuType,
	}
}

func (g *GPU) ID() types.GPUID {
	return g.gpuID
}

func (g *GPU) Type() types.GPUType {
	return g.gpuType
}

func (g *GPU) String() string {
	return fmt.Sprintf("gpu=[ID=%d, Type=%s]", g.ID(), g.Type())
}

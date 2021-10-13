package ue_procedures

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
)

func  NewIKESecurityAssociation() (uint64, error) {
	var maxSPI *big.Int = new(big.Int).SetUint64(math.MaxUint64)
	var localSPIuint64 uint64

	localSPI, err := rand.Int(rand.Reader, maxSPI)
	if err != nil {
		log.Fatalf("[Context] Error occurs when generate new IKE SPI")
		return 0, errors.New("Error occurs when generate new IKE SPI")
	}
	localSPIuint64 = localSPI.Uint64()

	fmt.Printf("\n\n***************************** SPI UE %d ******************************************** \n\n", localSPIuint64 )

	return localSPIuint64, nil
}

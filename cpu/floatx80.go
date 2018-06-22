package cpu

import (
	"math"
	"math/bits"
)

const (
	floatx80RoundingPrecision = 80
	floatx80Sign              = 1 << 15
)

const (
	// Software IEC/IEEE floating-point rounding mode.
	floatRoundNearestEven = iota
	floatRoundtoZero   
	floatRounddown       
	floatRoundup       
)

const (
	floatFlaginvalid = 1<<iota
	floatFlagDenormal 
	floatFlagDivbyzero 
	 floatFlagOverflow 
	 floatFlagUnderflow 
	  floatFlagInexact 
)

type (
	// floatx80 is a 80bit precision float
	floatx80 struct {
		high uint16
		low  uint64
	}
)

var (
	floatx80NaN = floatx80{0x7FFF, 1}
	floatRoundingMode = floatRoundNearestEven
	floatExceptionFlags = floatFlaginvalid
)

// *********************** converters **********************************

func int32ToFloatx80(a int32) floatx80 {
	if a == 0 {
		return packFloatx80(false, 0, 0)
	}
	zSign := a < 0
	if zSign {
		a = -a
	}
	shiftCount := uint16(bits.LeadingZeros32(uint32(a)) + 32)
	return packFloatx80(zSign, 0x403E-shiftCount,
		uint64(a)<<shiftCount)
}

func float32ToFloatx80(a float32) floatx80 {
	b := math.Float32bits(a)
	aFrac := uint64(0x007fffff & b)
	aExp := uint16(0xff & (b >> 23))
	aSign := (b & 0x80000000) != 0

	if aExp == 0xFF {
		return floatx80NaN
	}
	if aExp == 0 {
		if aFrac == 0 {
			return packFloatx80(aSign, 0, 0)
		}
		shiftCount := uint16(bits.LeadingZeros64(aFrac) - 32 - 8)
		aFrac <<= shiftCount
		aExp = 1 - shiftCount
	}
	aFrac |= 0x00800000
	return packFloatx80(aSign, aExp+0x3F80, aFrac<<40)
}

func float64ToFloatx80(a float64) floatx80 {
	b := math.Float64bits(a)
	aFrac := uint64(0x000FFFFFFFFFFFFF & b)
	aExp := uint16(0x7ff & (b >> 52))
	aSign := (b & 0x8000000000000000) != 0

	if aExp == 0x7FF {
		return floatx80NaN
	}
	if aExp == 0 {
		if aFrac == 0 {
			return packFloatx80(aSign, 0, 0)
		}
		shiftCount := uint16(bits.LeadingZeros64(aFrac) - 11)
		aFrac <<= shiftCount
		aExp = 1 - shiftCount
	}
	aFrac |= 0x0010000000000000
	return packFloatx80(aSign, aExp+0x3C00, aFrac<<11)
}

func (a floatx80) toInt32() int32 {
	aFrac := a.low
	aExp := a.high & 0x7fff
	aSign := (a.high & 0x8000) != 0
	if aExp == 0x07fff && (aFrac&0x8000000000000000) != 0 {
		aSign = false
	}
	shiftCount := 0x4037 - aExp
	if shiftCount <= 0 {
		shiftCount = 1
	}
	return roundAndPackInt32(aSign, shift64RightJamming(aFrac, shiftCount))
}

func (a floatx80) toInt32RoundToZero() int32 {
	return 0
}

func floatx80ToInt64(a floatx80) int64 {
	return 0
}

func floatx80ToInt64RoundToZero(a floatx80) int64 {
	return 0
}

func floatx80ToFloat32(a floatx80) float32 {
	return 0
}

func floatx80ToFloat32RoundToZero(a floatx80) float32 {
	return 0
}

func floatx80ToFloat64(a floatx80) float64 {
	return 0
}

func floatx80ToFloat64RoundToZero(a floatx80) float64 {
	return 0
}

/*
floatx80 floatx80_round_to_int( floatx80 );
floatx80 floatx80_add( floatx80, floatx80 );
floatx80 floatx80_sub( floatx80, floatx80 );
floatx80 floatx80_mul( floatx80, floatx80 );
floatx80 floatx80_div( floatx80, floatx80 );
floatx80 floatx80_rem( floatx80, floatx80 );
floatx80 floatx80_sqrt( floatx80 );
flag floatx80_eq( floatx80, floatx80 );
flag floatx80_le( floatx80, floatx80 );
flag floatx80_lt( floatx80, floatx80 );
flag floatx80_eq_signaling( floatx80, floatx80 );
flag floatx80_le_quiet( floatx80, floatx80 );
flag floatx80_lt_quiet( floatx80, floatx80 );
flag floatx80_is_signaling_nan( floatx80 );

int floatx80_fsin(floatx80 &a);
int floatx80_fcos(floatx80 &a);
int floatx80_ftan(floatx80 &a);

floatx80 floatx80_flognp1(floatx80 a);
floatx80 floatx80_flogn(floatx80 a);
floatx80 floatx80_flog2(floatx80 a);
floatx80 floatx80_flog10(floatx80 a);
*/

// *********************** helpers **********************************

func packFloatx80(zSign bool, zExp uint16, zSig uint64) floatx80 {
	if zSign {
		zExp += floatx80Sign
	}
	return floatx80{zExp, zSig}
}

func shift64RightJamming(a uint64, count uint16) uint64 {
	var z uint64

	if count == 0 {
		z = a
	} else if count < 64 {
		z = a >> count
		if (a << ((-count) & 63)) != 0 {
			z |= 1
		}
	} else {
		if a != 0 {
			z = 1
		} else {
			z = 0
		}
	}
	return z
}

func roundAndPackInt32(  zSign bool, absZ uint64) int32 {
	var z int32

	roundingMode := floatRoundingMode
	roundNearestEven := roundingMode == floatRoundNearestEven;
	roundIncrement := uint64(0x40);
	if ! roundNearestEven {
		if roundingMode == floatRoundtoZero {
			roundIncrement = 0;
		}  else {
			roundIncrement = 0x7F;
			if ( zSign ) {
				if roundingMode == floatRoundup {
					roundIncrement = 0
				}
			} else {
				if roundingMode == floatRounddown{
					  roundIncrement = 0
				}
			}
		}
	}
	roundBits := absZ & 0x7F
	absZ = ( absZ + roundIncrement )>>7
	absZ &= ~ ( ( ( roundBits ^ 0x40 ) == 0 ) & roundNearestEven )
	z = int32(absZ)
	if  zSign {  
		z = - z
	}
	if  absZ>>32 || z && ( ( z < 0 ) ^ zSign )  {
		floatRaise( floatFlaginvalid );
		if zSign {
			return -2147483648
		} else {
			return 0x7FFFFFFF
		}
	}
	if  roundBits!=0 {
		floatExceptionFlags |= floatFlagInexact
	}
	return z;
}

func floatRaise( floatFlag int ) {

}
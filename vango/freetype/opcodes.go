// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

// The Truetype opcodes are summarized at
// https://developer.apple.com/fonts/TTRefMan/RM07/appendixA.html

// TODO: the kOpXXX constants without end-of-line comments are not yet
// implemented.

const (
	kOpSVTCA0    = 0x00 // Set freedom and projection Vectors To Coordinate Axis
	kOpSVTCA1    = 0x01 // .
	kOpSPVTCA0   = 0x02 // Set Projection Vector To Coordinate Axis
	kOpSPVTCA1   = 0x03 // .
	kOpSFVTCA0   = 0x04 // Set Freedom Vector to Coordinate Axis
	kOpSFVTCA1   = 0x05 // .
	kOpSPVTL0    = 0x06 // Set Projection Vector To Line
	kOpSPVTL1    = 0x07 // .
	kOpSFVTL0    = 0x08 // Set Freedom Vector To Line
	kOpSFVTL1    = 0x09 // .
	kOpSPVFS     = 0x0a // Set Projection Vector From Stack
	kOpSFVFS     = 0x0b // Set Freedom Vector From Stack
	kOpGPV       = 0x0c // Get Projection Vector
	kOpGFV       = 0x0d // Get Freedom Vector
	kOpSFVTPV    = 0x0e // Set Freedom Vector To Projection Vector
	kOpISECT     = 0x0f
	kOpSRP0      = 0x10 // Set Reference Point 0
	kOpSRP1      = 0x11 // Set Reference Point 1
	kOpSRP2      = 0x12 // Set Reference Point 2
	kOpSZP0      = 0x13 // Set Zone Pointer 0
	kOpSZP1      = 0x14 // Set Zone Pointer 1
	kOpSZP2      = 0x15 // Set Zone Pointer 2
	kOpSZPS      = 0x16 // Set Zone PointerS
	kOpSLOOP     = 0x17 // Set LOOP variable
	kOpRTG       = 0x18 // Round To Grid
	kOpRTHG      = 0x19 // Round To Half Grid
	kOpSMD       = 0x1a // Set Minimum Distance
	kOpELSE      = 0x1b // ELSE clause
	kOpJMPR      = 0x1c // JuMP Relative
	kOpSCVTCI    = 0x1d // Set Control Value Table Cut-In
	kOpSSWCI     = 0x1e // Set Single Width Cut-In
	kOpSSW       = 0x1f // Set Single Width
	kOpDUP       = 0x20 // DUPlicate top stack element
	kOpPOP       = 0x21 // POP top stack element
	kOpCLEAR     = 0x22 // CLEAR the stack
	kOpSWAP      = 0x23 // SWAP the top two elements on the stack
	kOpDEPTH     = 0x24 // DEPTH of the stack
	kOpCINDEX    = 0x25 // Copy the INDEXed element to the top of the stack
	kOpMINDEX    = 0x26 // Move the INDEXed element to the top of the stack
	kOpALIGNPTS  = 0x27
	kOp_0x28     = 0x28
	kOpUTP       = 0x29
	kOpLOOPCALL  = 0x2a // LOOP and CALL function
	kOpCALL      = 0x2b // CALL function
	kOpFDEF      = 0x2c // Function DEFinition
	kOpENDF      = 0x2d // END Function definition
	kOpMDAP0     = 0x2e // Move Direct Absolute Point
	kOpMDAP1     = 0x2f // .
	kOpIUP0      = 0x30 // Interpolate Untouched Points through the outline
	kOpIUP1      = 0x31 // .
	kOpSHP0      = 0x32 // SHift Point using reference point
	kOpSHP1      = 0x33 // .
	kOpSHC0      = 0x34
	kOpSHC1      = 0x35
	kOpSHZ0      = 0x36 // SHift Zone using reference point
	kOpSHZ1      = 0x37 // .
	kOpSHPIX     = 0x38 // SHift point by a PIXel amount
	kOpIP        = 0x39 // Interpolate Point
	kOpMSIRP0    = 0x3a // Move Stack Indirect Relative Point
	kOpMSIRP1    = 0x3b // .
	kOpALIGNRP   = 0x3c // ALIGN to Reference Point
	kOpRTDG      = 0x3d // Round To Double Grid
	kOpMIAP0     = 0x3e // Move Indirect Absolute Point
	kOpMIAP1     = 0x3f // .
	kOpNPUSHB    = 0x40 // PUSH N Bytes
	kOpNPUSHW    = 0x41 // PUSH N Words
	kOpWS        = 0x42 // Write Store
	kOpRS        = 0x43 // Read Store
	kOpWCVTP     = 0x44 // Write Control Value Table in Pixel units
	kOpRCVT      = 0x45 // Read Control Value Table entry
	kOpGC0       = 0x46 // Get Coordinate projected onto the projection vector
	kOpGC1       = 0x47 // .
	kOpSCFS      = 0x48
	kOpMD0       = 0x49 // Measure Distance
	kOpMD1       = 0x4a // .
	kOpMPPEM     = 0x4b // Measure Pixels Per EM
	kOpMPS       = 0x4c // Measure Point Size
	kOpFLIPON    = 0x4d // set the auto FLIP Boolean to ON
	kOpFLIPOFF   = 0x4e // set the auto FLIP Boolean to OFF
	kOpDEBUG     = 0x4f // DEBUG call
	kOpLT        = 0x50 // Less Than
	kOpLTEQ      = 0x51 // Less Than or EQual
	kOpGT        = 0x52 // Greater Than
	kOpGTEQ      = 0x53 // Greater Than or EQual
	kOpEQ        = 0x54 // EQual
	kOpNEQ       = 0x55 // Not EQual
	kOpODD       = 0x56 // ODD
	kOpEVEN      = 0x57 // EVEN
	kOpIF        = 0x58 // IF test
	kOpEIF       = 0x59 // End IF
	kOpAND       = 0x5a // logical AND
	kOpOR        = 0x5b // logical OR
	kOpNOT       = 0x5c // logical NOT
	kOpDELTAP1   = 0x5d // DELTA exception P1
	kOpSDB       = 0x5e // Set Delta Base in the graphics state
	kOpSDS       = 0x5f // Set Delta Shift in the graphics state
	kOpADD       = 0x60 // ADD
	kOpSUB       = 0x61 // SUBtract
	kOpDIV       = 0x62 // DIVide
	kOpMUL       = 0x63 // MULtiply
	kOpABS       = 0x64 // ABSolute value
	kOpNEG       = 0x65 // NEGate
	kOpFLOOR     = 0x66 // FLOOR
	kOpCEILING   = 0x67 // CEILING
	kOpROUND00   = 0x68 // ROUND value
	kOpROUND01   = 0x69 // .
	kOpROUND10   = 0x6a // .
	kOpROUND11   = 0x6b // .
	kOpNROUND00  = 0x6c // No ROUNDing of value
	kOpNROUND01  = 0x6d // .
	kOpNROUND10  = 0x6e // .
	kOpNROUND11  = 0x6f // .
	kOpWCVTF     = 0x70
	kOpDELTAP2   = 0x71 // DELTA exception P2
	kOpDELTAP3   = 0x72 // DELTA exception P3
	kOpDELTAC1   = 0x73
	kOpDELTAC2   = 0x74
	kOpDELTAC3   = 0x75
	kOpSROUND    = 0x76 // Super ROUND
	kOpS45ROUND  = 0x77 // Super ROUND 45 degrees
	kOpJROT      = 0x78 // Jump Relative On True
	kOpJROF      = 0x79 // Jump Relative On False
	kOpROFF      = 0x7a // Round OFF
	kOp_0x7b     = 0x7b
	kOpRUTG      = 0x7c // Round Up To Grid
	kOpRDTG      = 0x7d // Round Down To Grid
	kOpSANGW     = 0x7e // Set ANGle Weight
	kOpAA        = 0x7f // Adjust Angle
	kOpFLIPPT    = 0x80
	kOpFLIPRGON  = 0x81
	kOpFLIPRGOFF = 0x82
	kOp_0x83     = 0x83
	kOp_0x84     = 0x84
	kOpSCANCTRL  = 0x85 // SCAN conversion ConTRoL
	kOpSDPVTL0   = 0x86
	kOpSDPVTL1   = 0x87
	kOpGETINFO   = 0x88 // GET INFOrmation
	kOpIDEF      = 0x89 // Instruction DEFinition
	kOpROLL      = 0x8a // ROLL the top three stack elements
	kOpMAX       = 0x8b // MAXimum of top two stack elements
	kOpMIN       = 0x8c // MINimum of top two stack elements
	kOpSCANTYPE  = 0x8d // SCANTYPE
	kOpINSTCTRL  = 0x8e
	kOp0x8f      = 0x8f
	kOp0x90      = 0x90
	kOp0x91      = 0x91
	kOp0x92      = 0x92
	kOp0x93      = 0x93
	kOp0x94      = 0x94
	kOp0x95      = 0x95
	kOp0x96      = 0x96
	kOp0x97      = 0x97
	kOp0x98      = 0x98
	kOp0x99      = 0x99
	kOp0x9a      = 0x9a
	kOp0x9b      = 0x9b
	kOp0x9c      = 0x9c
	kOp0x9d      = 0x9d
	kOp0x9e      = 0x9e
	kOp0x9f      = 0x9f
	kOp0xa0      = 0xa0
	kOp0xa1      = 0xa1
	kOp0xa2      = 0xa2
	kOp0xa3      = 0xa3
	kOp0xa4      = 0xa4
	kOp0xa5      = 0xa5
	kOp0xa6      = 0xa6
	kOp0xa7      = 0xa7
	kOp0xa8      = 0xa8
	kOp0xa9      = 0xa9
	kOp0xaa      = 0xaa
	kOp0xab      = 0xab
	kOp0xac      = 0xac
	kOp0xad      = 0xad
	kOp0xae      = 0xae
	kOp0xaf      = 0xaf
	kOpPUSHB000  = 0xb0 // PUSH Bytes
	kOpPUSHB001  = 0xb1 // .
	kOpPUSHB010  = 0xb2 // .
	kOpPUSHB011  = 0xb3 // .
	kOpPUSHB100  = 0xb4 // .
	kOpPUSHB101  = 0xb5 // .
	kOpPUSHB110  = 0xb6 // .
	kOpPUSHB111  = 0xb7 // .
	kOpPUSHW000  = 0xb8 // PUSH Words
	kOpPUSHW001  = 0xb9 // .
	kOpPUSHW010  = 0xba // .
	kOpPUSHW011  = 0xbb // .
	kOpPUSHW100  = 0xbc // .
	kOpPUSHW101  = 0xbd // .
	kOpPUSHW110  = 0xbe // .
	kOpPUSHW111  = 0xbf // .
	kOpMDRP00000 = 0xc0 // Move Direct Relative Point
	kOpMDRP00001 = 0xc1 // .
	kOpMDRP00010 = 0xc2 // .
	kOpMDRP00011 = 0xc3 // .
	kOpMDRP00100 = 0xc4 // .
	kOpMDRP00101 = 0xc5 // .
	kOpMDRP00110 = 0xc6 // .
	kOpMDRP00111 = 0xc7 // .
	kOpMDRP01000 = 0xc8 // .
	kOpMDRP01001 = 0xc9 // .
	kOpMDRP01010 = 0xca // .
	kOpMDRP01011 = 0xcb // .
	kOpMDRP01100 = 0xcc // .
	kOpMDRP01101 = 0xcd // .
	kOpMDRP01110 = 0xce // .
	kOpMDRP01111 = 0xcf // .
	kOpMDRP10000 = 0xd0 // .
	kOpMDRP10001 = 0xd1 // .
	kOpMDRP10010 = 0xd2 // .
	kOpMDRP10011 = 0xd3 // .
	kOpMDRP10100 = 0xd4 // .
	kOpMDRP10101 = 0xd5 // .
	kOpMDRP10110 = 0xd6 // .
	kOpMDRP10111 = 0xd7 // .
	kOpMDRP11000 = 0xd8 // .
	kOpMDRP11001 = 0xd9 // .
	kOpMDRP11010 = 0xda // .
	kOpMDRP11011 = 0xdb // .
	kOpMDRP11100 = 0xdc // .
	kOpMDRP11101 = 0xdd // .
	kOpMDRP11110 = 0xde // .
	kOpMDRP11111 = 0xdf // .
	kOpMIRP00000 = 0xe0 // Move Indirect Relative Point
	kOpMIRP00001 = 0xe1 // .
	kOpMIRP00010 = 0xe2 // .
	kOpMIRP00011 = 0xe3 // .
	kOpMIRP00100 = 0xe4 // .
	kOpMIRP00101 = 0xe5 // .
	kOpMIRP00110 = 0xe6 // .
	kOpMIRP00111 = 0xe7 // .
	kOpMIRP01000 = 0xe8 // .
	kOpMIRP01001 = 0xe9 // .
	kOpMIRP01010 = 0xea // .
	kOpMIRP01011 = 0xeb // .
	kOpMIRP01100 = 0xec // .
	kOpMIRP01101 = 0xed // .
	kOpMIRP01110 = 0xee // .
	kOpMIRP01111 = 0xef // .
	kOpMIRP10000 = 0xf0 // .
	kOpMIRP10001 = 0xf1 // .
	kOpMIRP10010 = 0xf2 // .
	kOpMIRP10011 = 0xf3 // .
	kOpMIRP10100 = 0xf4 // .
	kOpMIRP10101 = 0xf5 // .
	kOpMIRP10110 = 0xf6 // .
	kOpMIRP10111 = 0xf7 // .
	kOpMIRP11000 = 0xf8 // .
	kOpMIRP11001 = 0xf9 // .
	kOpMIRP11010 = 0xfa // .
	kOpMIRP11011 = 0xfb // .
	kOpMIRP11100 = 0xfc // .
	kOpMIRP11101 = 0xfd // .
	kOpMIRP11110 = 0xfe // .
	kOpMIRP11111 = 0xff // .
)

// kPopNumOfStack[opcode] == ff means that that opcode is not yet implemented.
const ff = 255
const kOpN = ff

// kPopNum is the number of stack elements that each opcode pops.
var kPopNum = [256]uint8{
	// 1, 2, 3, 4, 5, 6, 7, 8, 9, a, b, c, d, e, f
	00, 00, 00, 00, 00, 00, 02, 02, 02, 02, 02, 02, 00, 00, 00, ff, // 0x00-0x0f
	01, 01, 01, 01, 01, 01, 01, 01, 00, 00, 01, 00, 01, 01, 01, 01, // 0x10-0x1f
	01, 01, 00, 02, 00, 01, 01, ff, ff, ff, 02, 01, 01, 00, 01, 01, // 0x20-0x2f
	00, 00, 00, 00, ff, ff, 01, 01, 01, 00, 02, 02, 00, 00, 02, 02, // 0x30-0x3f
	00, 00, 02, 01, 02, 01, 01, 01, ff, 02, 02, 00, 00, 00, 00, 00, // 0x40-0x4f
	02, 02, 02, 02, 02, 02, 01, 01, 01, 00, 02, 02, 01, 01, 01, 01, // 0x50-0x5f
	02, 02, 02, 02, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, // 0x60-0x6f
	ff, 01, 01, ff, ff, ff, 01, 01, 02, 02, 00, ff, 00, 00, 01, 01, // 0x70-0x7f
	ff, ff, ff, ff, ff, 01, ff, ff, 01, 01, 03, 02, 02, 01, ff, ff, // 0x80-0x8f
	ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, // 0x90-0x9f
	ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, ff, // 0xa0-0xaf
	00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, // 0xb0-0xbf
	01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, // 0xc0-0xcf
	01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, 01, // 0xd0-0xdf
	02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, // 0xe0-0xef
	02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, 02, // 0xf0-0xff
}

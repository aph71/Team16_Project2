package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	rawInstruction  string
	lineValue       uint64
	memLoc          uint64
	opcode          uint64
	op              string
	instructionType string
	rm              uint8
	shamt           uint8
	rn              uint8
	rd              uint8
	rt              uint8
	op2             uint8
	address         uint16
	immediate       int16
	offset          int32
	conditional     uint8
	shiftCode       uint8
	field           uint32
	memValue        int64
}

func ReadBinary(filePath string) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var linenumber uint64
	linenumber = 96
	for scanner.Scan() {
		instruction := strings.ReplaceAll(scanner.Text(), " ", "")
		InstructionList = append(InstructionList, Instruction{rawInstruction: instruction, memLoc: linenumber})
		linenumber += 4
	}
}

func WriteInstructions(filePath string, list []Instruction) {
	outputFile, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer outputFile.Close()

	for i := 0; i < len(list); i++ {
		switch list[i].instructionType {
		case "B":
			//write binary with spaces
			_, err := fmt.Fprintf(outputFile, "%s %s\t", list[i].rawInstruction[0:6], list[i].rawInstruction[6:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(outputFile, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(outputFile, "#%d\n", list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "I":
			//write binary with spaces
			_, err := fmt.Fprintf(outputFile, "%s %s %s %s\t", list[i].rawInstruction[0:10], list[i].rawInstruction[10:22], list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(outputFile, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(outputFile, "R%d, R%d, #%d\n", list[i].rd, list[i].rn, list[i].immediate)
			if err != nil {
				log.Fatal(err)
			}

		case "CB":
			//write binary with spaces
			_, err := fmt.Fprintf(outputFile, "%s %s %s\t", list[i].rawInstruction[0:8], list[i].rawInstruction[8:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(outputFile, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(outputFile, "R%d, #%d\n", list[i].conditional, list[i].offset)
			if err != nil {
				log.Fatal(err)
			}
		case "IM":
			//write binary with spaces
			_, err := fmt.Fprintf(outputFile, "%s %s %s %s\t", list[i].rawInstruction[0:9], list[i].rawInstruction[9:12], list[i].rawInstruction[12:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(outputFile, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(outputFile, "R%d, %d, LSL %d\n", list[i].rd, list[i].field, list[i].shiftCode)
			if err != nil {
				log.Fatal(err)
			}
		case "D":
			//write binary with spaces
			_, err := fmt.Fprintf(outputFile, "%s %s %s %s %s\t", list[i].rawInstruction[0:11], list[i].rawInstruction[11:20], list[i].rawInstruction[20:22], list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(outputFile, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(outputFile, "R%d, [R%d, #%d]\n", list[i].rt, list[i].rn, list[i].address)
			if err != nil {
				log.Fatal(err)
			}
		case "R":
			//write binary with spaces
			_, err := fmt.Fprintf(outputFile, "%s %s %s %s %s\t", list[i].rawInstruction[0:11], list[i].rawInstruction[11:16], list[i].rawInstruction[16:22], list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			//write memLoc and opcode
			_, err = fmt.Fprintf(outputFile, "%d\t%s\t", list[i].memLoc, list[i].op)
			//write operands
			_, err = fmt.Fprintf(outputFile, "R%d, R%d, ", list[i].rd, list[i].rn)
			if list[i].op == "LSL" || list[i].op == "ASR" || list[i].op == "LSR" {
				_, err = fmt.Fprintf(outputFile, "#%d\n", list[i].shamt)
			} else {
				_, err = fmt.Fprintf(outputFile, "R%d\n", list[i].rm)
			}
			if err != nil {
				log.Fatal(err)
			}
		case "BREAK":
			_, err := fmt.Fprintf(outputFile, "%s\t%d\tBREAK\n", list[i].rawInstruction, list[i].memLoc)
			if err != nil {
				log.Fatal(err)
			}
		case "MEM":
			_, err := fmt.Fprintf(outputFile, "%s\t%d\t%d\n", list[i].rawInstruction, list[i].memLoc, list[i].memValue)
			if err != nil {
				log.Fatal(err)
			}
		case "NOP":
			_, err := fmt.Fprintf(outputFile, "%s\t%d\t%s\n", list[i].rawInstruction, list[i].memLoc, list[i].op)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func ProcessInstructionList(list []Instruction) {
	breakHit := false
	for i := 0; i < len(list); i++ {
		if !breakHit {
			translateToInt(&list[i])
			opcodeMasking(&list[i])
			opcodeTranslation(&list[i])
			switch list[i].instructionType {
			case "B":
				processBType(&list[i])
			case "I":
				processIType(&list[i])
			case "CB":
				processCBType(&list[i])
			case "IM":
				processIMType(&list[i])
			case "D":
				processDType(&list[i])
			case "R":
				processRType(&list[i])
			case "BREAK":
				breakHit = true
			}
		} else {
			list[i].instructionType = "MEM"
			var value uint64
			value, _ = strconv.ParseUint(list[i].rawInstruction, 2, 32)
			list[i].memValue = parse2Complement(value, 32)
		}
	}
}

// translates raw instruction to an unsigned 64 bit int
func translateToInt(ins *Instruction) {
	i, err := strconv.ParseUint(ins.rawInstruction, 2, 64)
	if err == nil {
		ins.lineValue = i
	} else {
		fmt.Println(err)
	}
}

func opcodeMasking(ins *Instruction) {
	ins.opcode = (ins.lineValue & 4292870144) >> 21
}

func opcodeTranslation(ins *Instruction) {
	if ins.opcode >= 160 && ins.opcode <= 191 {
		ins.op = "B"
		ins.instructionType = "B"
	} else if ins.opcode == 1104 {
		ins.op = "AND"
		ins.instructionType = "R"
	} else if ins.opcode == 1112 {
		ins.op = "ADD"
		ins.instructionType = "R"
	} else if ins.opcode >= 1160 && ins.opcode <= 1161 {
		ins.op = "ADDI"
		ins.instructionType = "I"
	} else if ins.opcode == 1360 {
		ins.op = "ORR"
		ins.instructionType = "R"
	} else if ins.opcode >= 1440 && ins.opcode <= 1447 {
		ins.op = "CBZ"
		ins.instructionType = "CB"
	} else if ins.opcode >= 1448 && ins.opcode <= 1455 {
		ins.op = "CBNZ"
		ins.instructionType = "CB"
	} else if ins.opcode == 1624 {
		ins.op = "SUB"
		ins.instructionType = "R"
	} else if ins.opcode >= 1672 && ins.opcode <= 1673 {
		ins.op = "SUBI"
		ins.instructionType = "I"
	} else if ins.opcode >= 1684 && ins.opcode <= 1687 {
		ins.op = "MOVZ"
		ins.instructionType = "IM"
	} else if ins.opcode >= 1940 && ins.opcode <= 1943 {
		ins.op = "MOVK"
		ins.instructionType = "IM"
	} else if ins.opcode == 1690 {
		ins.op = "LSR"
		ins.instructionType = "R"
	} else if ins.opcode == 1691 {
		ins.op = "LSL"
		ins.instructionType = "R"
	} else if ins.opcode == 1984 {
		ins.op = "STUR"
		ins.instructionType = "D"
	} else if ins.opcode == 1986 {
		ins.op = "LDUR"
		ins.instructionType = "D"
	} else if ins.opcode == 1692 {
		ins.op = "ASR"
		ins.instructionType = "R"
	} else if ins.opcode == 0 {
		ins.op = "NOP"
		ins.instructionType = "NOP"
	} else if ins.opcode == 1872 {
		ins.op = "EOR"
		ins.instructionType = "R"
	} else if ins.opcode == 2038 {
		ins.op = "BREAK"
		ins.instructionType = "BREAK"
	} else if ins.opcode == 0 {
		ins.op = "NOP"
		ins.instructionType = "NOP"
	} else {
		fmt.Println("Invalid opcode")
	}
}

func processRType(ins *Instruction) {
	//mask for bits 12 - 16
	ins.rm = uint8((ins.lineValue & 2031616) >> 16)
	//mask for bits 17 - 22
	ins.shamt = uint8((ins.lineValue & 64512) >> 10)
	//mask for bits 23 - 27
	ins.rn = uint8((ins.lineValue & 992) >> 5)
	//mask for bit 28 - 32
	ins.rd = uint8(ins.lineValue & 31)
}

func processIType(ins *Instruction) {
	//mask for bits 11 - 22
	ins.immediate = int16(parse2Complement((ins.lineValue&4193280)>>10, 12))
	//mask for bits 23 - 27
	ins.rn = uint8((ins.lineValue & 992) >> 5)
	//mask for bits 28 - 32
	ins.rd = uint8(ins.lineValue & 31)
}

func processCBType(ins *Instruction) {
	//mask for bits 9 - 27
	ins.offset = int32(parse2Complement((ins.lineValue&16777184)>>5, 19))
	//mask for bits 28 - 32
	ins.conditional = uint8(ins.lineValue & 31)
}

func processIMType(ins *Instruction) {
	//mask for bits 10 - 12
	ins.shiftCode = uint8((ins.lineValue & 6291456) >> 21)
	//mask for bits 13 - 27
	ins.field = uint32((ins.lineValue & 2097120) >> 5)
	//mask for bits 28 - 32
	ins.rd = uint8(ins.lineValue & 31)
}

func processDType(ins *Instruction) {
	//mask for bits 12 - 20
	ins.address = uint16((ins.lineValue & 2093056) >> 12)
	//mask for bits 21 - 22
	ins.op2 = uint8((ins.lineValue & 3072) >> 10)
	//mask for bits 23 - 27
	ins.rn = uint8((ins.lineValue & 992) >> 5)
	//mask for bit 28 - 32
	ins.rt = uint8(ins.lineValue & 31)
}

func processBType(ins *Instruction) {
	//mask for bits 7 - 32
	ins.offset = int32(parse2Complement(ins.lineValue&67108863, 26))
}

// parses 2's complement binary to an integer
func parse2Complement(i uint64, binaryLength uint) int64 {
	var out int64
	var xorValue int64
	out = int64(i)
	xorValue = (1 << binaryLength) - 1
	if (i >> (binaryLength - 1)) != 0 {
		out = ((out ^ xorValue) + 1) * -1
	}
	return out
}

func outputRegistersToFile(registry []int, simOutputFile *os.File, myMap map[int]map[int]int) {
	fmt.Fprintf(simOutputFile, "registers:\n")
	for i := 0; i < len(registry); i += 8 {
		end := i + 8
		if end > len(registry) {
			end = len(registry)
		}
		fmt.Fprintf(simOutputFile, "r%d:\t", i)
		for j := i; j < end; j++ {
			fmt.Fprintf(simOutputFile, "%d\t", registry[j])
		}
		fmt.Fprintf(simOutputFile, "\n")
	}
	if len(myMap) > 0 {
		fmt.Fprintf(simOutputFile, "data:\n")
		for k1, v1 := range myMap {
			fmt.Fprintf(simOutputFile, "%d:", k1)
			for k2 := 0; k2 < 8; k2++ { // Ensure iteration for keys 0 to 7
				val, exists := v1[k2]
				if exists {
					fmt.Fprintf(simOutputFile, "%d\t", val)
				} else {
					fmt.Fprintf(simOutputFile, "0\t")
				}
			}
			fmt.Fprintf(simOutputFile, "\n")
		}
	}
	fmt.Fprintf(simOutputFile, "============\n")
}

func simulateInstruction(simOutput string, list []Instruction, registry []int, myMap map[int]map[int]int) {
	breakHit := false
	cycle := 1
	simOutputFile, err := os.Create(simOutput)

	if err != nil {
		log.Fatal(err)
	}

	defer simOutputFile.Close()
	for i := 0; i < len(list); i++ {
		if !breakHit {
			switch opcode := list[i].opcode; {
			//*****B INSTRUCTION****
			case opcode >= 160 && opcode <= 191:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t#%d\n", cycle, list[i].memLoc, list[i].op, list[i].offset)
				outputRegistersToFile(registry, simOutputFile, otherData)
				cycle++
				i = i + int(list[i].offset-1)
			//*****AND INSTRUCTION*****
			case opcode == 1104:
				regDest := registry[list[i].rn] & registry[list[i].rm]
				registry[list[i].rd] = regDest
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s R%d, R%d, R%d\n", cycle, list[i].memLoc, list[i].op, list[i].rd, list[i].rn, list[i].rm)
				outputRegistersToFile(registry, simOutputFile, otherData)

			//*****ADD INSTRUCTION*****
			case opcode == 1112:
				//fmt.Println(list[i].rn)
				//ADDED FOR TESTING PURPOSES ONLY
				//END TESTING BLOCK
				regDest := registry[list[i].rm] + registry[list[i].rn]
				registry[list[i].rd] = regDest
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s R%d, R%d, R%d\n", cycle, list[i].memLoc, list[i].op, list[i].rd, list[i].rn, list[i].rm)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****ADDI INSTRUCTION*****
			case opcode == 1160 || opcode == 1161:
				regDest := registry[list[i].rn] + int(list[i].immediate)
				registry[list[i].rd] = regDest
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s R%d, R%d, R%d\n", cycle, list[i].memLoc, list[i].op, list[i].rd, list[i].rn, list[i].immediate)
				outputRegistersToFile(registry, simOutputFile, otherData)
				cycle++
				//*****ORR INSTRUCTION*****
			case opcode == 1360:
				regDest := registry[list[i].rn] | registry[list[i].rm]
				registry[list[i].rd] = regDest
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s R%d, R%d, R%d\n", cycle, list[i].memLoc, list[i].op, list[i].rd, list[i].rn, list[i].rm)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****CBZ INSTRUCTION*****
			case opcode >= 1440 && opcode <= 1447:
				{
					fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\tR%d, #%d\n", cycle, list[i].memLoc, list[i].op, list[i].conditional, list[i].offset)
					outputRegistersToFile(registry, simOutputFile, otherData)

					cycle++
					if registry[list[i].conditional] == 0 {
						i = i + int(list[i].offset-1)
					}
				}
				//*****CBNZ*****
			case opcode >= 1448 && opcode <= 1455:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\tR%d, #%d\n", cycle, list[i].memLoc, list[i].op, list[i].conditional, list[i].offset)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				if registry[list[i].conditional] != 0 {
					i = i + int(list[i].offset-1)
				}
				//*****SUB INSTRUCTION*****
			case opcode == 1624:
				regDest := registry[list[i].rm] - registry[list[i].rn]
				registry[list[i].rd] = regDest
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****SUBI INSTRUCTION*****
			case opcode == 1672 || opcode == 1673:
				regDest := registry[list[i].rn] - int(list[i].immediate)
				registry[list[i].rd] = regDest
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s R%d, R%d, #%d\n", cycle, list[i].memLoc, list[i].op, list[i].rd, list[i].rn, list[i].immediate)
				outputRegistersToFile(registry, simOutputFile, otherData)
				cycle++
				//*****MOVZ*****
			case opcode >= 1684 && opcode <= 1687:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****MOVK*****
			case opcode >= 1940 && opcode <= 1943:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****LSR*****
			case opcode == 1690:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****LSL*****
			case opcode == 1691:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****STUR*****
			case opcode == 1984:
				destAddress := registry[list[i].rn] + int((list[i].address)*4)
				data := registry[list[i].rt]
				if _, ok := myMap[destAddress]; !ok {
					myMap[destAddress] = make(map[int]int)
					for i := 0; i < 8; i++ {
						myMap[destAddress][i] = 0
						myMap[destAddress][0] = data
					}
				}
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\tR%d, [R%d, #%d]\n", cycle, list[i].memLoc, list[i].op, list[i].rt, list[i].rn,
					list[i].address)
				outputRegistersToFile(registry, simOutputFile, otherData)
				cycle++
				//*****LDUR*****
			case opcode == 1986:
				destAddress := registry[list[i].rn] + int((list[i].address)*4)
				if _, ok := myMap[destAddress]; !ok {
					// If the key doesn't exist, create the inner map and initialize 8 values to 0
					data := 0
					registry[list[i].rt] = data
				} else {
					data := myMap[destAddress][0]
					registry[list[i].rt] = data
				}
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\tR%d, [R%d, #%d]\n", cycle, list[i].memLoc, list[i].op, list[i].rt,
					list[i].rn, list[i].address)
				outputRegistersToFile(registry, simOutputFile, otherData)
				cycle++
				//*****ASR*****
			case opcode == 1692:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)
				cycle++
				//*****NOP*****
			case opcode == 0:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)
				cycle++
				//*****EOR*****
			case opcode == 1872:
				regDest := registry[list[i].rn] ^ registry[list[i].rm]
				registry[list[i].rd] = regDest
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s R%d, R%d, R%d\n", cycle, list[i].memLoc, list[i].op, list[i].rd, list[i].rn, list[i].rm)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++
				//*****BREAK*****
			case opcode == 2038:
				fmt.Fprintf(simOutputFile, "Cycle:%d\t%d\t%s\t\n", cycle, list[i].memLoc, list[i].op)
				outputRegistersToFile(registry, simOutputFile, otherData)

				cycle++

			}
		}
	}
}

var InstructionList []Instruction
var registryData = make([]int, 32)
var otherData = make(map[int]map[int]int)

func main() {
	inputFilePathPtr := flag.String("i", "input.txt", "input file path")
	outputFilePathPtr := flag.String("o", "out.txt", "output file path")
	outputFile2PathPtr := flag.String("o2", "outputsim.txt", "output sim file path")
	flag.Parse()
	ReadBinary(*inputFilePathPtr)
	ProcessInstructionList(InstructionList)
	//fmt.Println(InstructionList)
	WriteInstructions(*outputFilePathPtr, InstructionList)
	simulateInstruction(*outputFile2PathPtr, InstructionList, registryData, otherData)
}

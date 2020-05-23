package logic

import (
	"fmt"

	pb "github.com/ememak/Projekt-Rada/query"
)

// AcceptToken checks if MessageToSign is matching server informations.
//
// Function returns true if token is a token of data[qNum] (provided
// such query exists). In other case it returns false.
func AcceptToken(token *pb.VoteToken, qNum int32, data []pb.PollQuestion) bool {
	// Check if query of number qNum exists.
	if qNum >= int32(len(data)) || qNum < 0 {
		fmt.Printf("No such Query: %v\n", qNum)
		return false
	}

	// Search if data contains token.
	// For now data is just an array, this will be changed.
	nT := len(data[qNum].Tokens)
	if nT == 0 {
		fmt.Printf("Token not valid\n")
		return false
	}
	for i := 0; i < nT; i++ {
		if data[qNum].Tokens[i].Token == token.Token {
			data[qNum].Tokens = append(data[qNum].Tokens[:i], data[qNum].Tokens[i+1:]...)
			break
		}
		if i == nT-1 {
			fmt.Printf("Token not valid\n")
			return false
		}
	}
	return true
}

func AcceptVote(v *pb.Vote, data []pb.PollQuestion) (*pb.VoteReply, error) {
	if v.Nr >= int32(len(data)) || v.Nr < 0 { //security leak, path that out later!
		fmt.Printf("No such Query: %v\n", v.Nr)
		return &pb.VoteReply{Mess: "No such Query!\n"}, nil
	}
	var nF = len(data[v.Nr].Fields)
	if nF != len(v.Answer) {
		fmt.Printf("Vote have different number of fields than query\n")
		return &pb.VoteReply{Mess: "Vote have different number of fields than query!\n"}, nil
	}
	for i := 0; i < nF; i++ {
		if v.Answer[i] >= 1 {
			data[v.Nr].Fields[i].Votes++
		}
	}
	fmt.Printf("Thank you for your vote!\n")
	fmt.Printf("In Memory: %v\n", data)
	return &pb.VoteReply{Mess: "Thank you for your vote!\n"}, nil
}

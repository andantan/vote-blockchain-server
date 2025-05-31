genProtoVoteSubmit:
	@protoc vote_submit_message.proto --go_out=. --go-grpc_out=.

genProtoVoteProposal:
	@protoc vote_proposal_message.proto --go_out=. --go-grpc_out=.

genProto:
	genProtoVoteSubmit
	genProtoVoteProposal
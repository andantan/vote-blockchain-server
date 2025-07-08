.PHONY: genProtoVoteSubmit genProtoVoteProposal genProtoAdminCommand genProto

genProtoVoteSubmit:
	@protoc vote_submit_message.proto --go_out=. --go-grpc_out=.

genProtoVoteProposal:
	@protoc vote_proposal_message.proto --go_out=. --go-grpc_out=.

genProtoAdminCommand:
	@protoc admin_l4_commands.proto --go_out=. --go-grpc_out=.

genProto: genProtoVoteSubmit genProtoVoteProposal genProtoAdminCommand
	@echo "All proto files generated."
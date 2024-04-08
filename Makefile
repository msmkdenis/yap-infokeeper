.PHONY: all

# Generate mock service and repository for user
gen-mock-user-repo:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/user/mocks/mock_user_repository.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/user/service UserRepository

gen-mock-user-service:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/user/mocks/mock_user_service.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/user/api/grpchandlers UserService

# Generate mock service and repository for text data
gen-mock-text-data-service:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/text_data/mocks/mock_text_data_service.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers TextDataService

gen-mock-text-data-repo:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/text_data/mocks/mock_text_data_repository.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/text_data/service TextDataRepository

# Generate mock service and repository for credit card
gen-mock-credit-card-service:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/credit_card/mocks/mock_credit_card_service.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers CreditCardService

gen-mock-credit-card-repo:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/credit_card/mocks/mock_credit_card_repository.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credit_card/service CreditCardRepository

# Generate mock service and repository for credential
gen-mock-credential-repo:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/credential/mocks/mock_credential_repository.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credential/service CredentialRepository

gen-mock-credential-service:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/credential/mocks/mock_credential_service.go \
			 -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers CredentialService

# Generate proto for credit card
gen-proto-credit-card:
	@protoc --go_out=. --go_opt=paths=source_relative \
       		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
       		internal/credit_card/api/grpchandlers/proto/credit_card.proto

# Generate proto for user
gen-proto-user:
	@protoc --go_out=. --go_opt=paths=source_relative \
       		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
       		internal/user/api/grpchandlers/proto/user.proto

# Generate proto for text data
gen-proto-text-data:
	@protoc --go_out=. --go_opt=paths=source_relative \
       		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
       		internal/text_data/api/grpchandlers/proto/text_data.proto

# Generate proto for credential
gen-proto-credential:
	@protoc --go_out=. --go_opt=paths=source_relative \
       		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
       		internal/credential/api/grpchandlers/proto/credential.proto
protoc --proto_path=proto/ --go_out=. chat.proto
protoc --proto_path=proto/ --js_out="library=chat_pb,binary:chat" proto/chat.proto
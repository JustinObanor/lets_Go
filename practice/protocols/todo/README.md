They help us seriealize(pulbish, broadcast) data - like XML, but way better

protocol buffers is a format to store messages. sth like json or xml
message type and values of that message type
they have schema and value
file defines how it looks like. value of those schema

protoc --go_out=. todo.proto

go install ./cmd/todo  
#to install

todo
#to run

todo +flag
#to run with flag

todo add <text>

cat mydb.pb | protoc --decode_raw
#see whats inside a protocobuf
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: workerpb/worker.proto

package workerpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TaskRequest struct {
	Node                 string   `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Condition            string   `protobuf:"bytes,2,opt,name=condition,proto3" json:"condition,omitempty"`
	Action               string   `protobuf:"bytes,3,opt,name=action,proto3" json:"action,omitempty"`
	Params               string   `protobuf:"bytes,4,opt,name=params,proto3" json:"params,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskRequest) Reset()         { *m = TaskRequest{} }
func (m *TaskRequest) String() string { return proto.CompactTextString(m) }
func (*TaskRequest) ProtoMessage()    {}
func (*TaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e70f8309ff9e66, []int{0}
}

func (m *TaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskRequest.Unmarshal(m, b)
}
func (m *TaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskRequest.Marshal(b, m, deterministic)
}
func (m *TaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskRequest.Merge(m, src)
}
func (m *TaskRequest) XXX_Size() int {
	return xxx_messageInfo_TaskRequest.Size(m)
}
func (m *TaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TaskRequest proto.InternalMessageInfo

func (m *TaskRequest) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *TaskRequest) GetCondition() string {
	if m != nil {
		return m.Condition
	}
	return ""
}

func (m *TaskRequest) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *TaskRequest) GetParams() string {
	if m != nil {
		return m.Params
	}
	return ""
}

type TaskResult struct {
	Node                 string               `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Condition            string               `protobuf:"bytes,2,opt,name=condition,proto3" json:"condition,omitempty"`
	Action               string               `protobuf:"bytes,3,opt,name=action,proto3" json:"action,omitempty"`
	Worker               string               `protobuf:"bytes,4,opt,name=worker,proto3" json:"worker,omitempty"`
	Success              bool                 `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TaskResult) Reset()         { *m = TaskResult{} }
func (m *TaskResult) String() string { return proto.CompactTextString(m) }
func (*TaskResult) ProtoMessage()    {}
func (*TaskResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e70f8309ff9e66, []int{1}
}

func (m *TaskResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskResult.Unmarshal(m, b)
}
func (m *TaskResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskResult.Marshal(b, m, deterministic)
}
func (m *TaskResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskResult.Merge(m, src)
}
func (m *TaskResult) XXX_Size() int {
	return xxx_messageInfo_TaskResult.Size(m)
}
func (m *TaskResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskResult.DiscardUnknown(m)
}

var xxx_messageInfo_TaskResult proto.InternalMessageInfo

func (m *TaskResult) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *TaskResult) GetCondition() string {
	if m != nil {
		return m.Condition
	}
	return ""
}

func (m *TaskResult) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *TaskResult) GetWorker() string {
	if m != nil {
		return m.Worker
	}
	return ""
}

func (m *TaskResult) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *TaskResult) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func init() {
	proto.RegisterType((*TaskRequest)(nil), "worker.TaskRequest")
	proto.RegisterType((*TaskResult)(nil), "worker.TaskResult")
}

func init() { proto.RegisterFile("workerpb/worker.proto", fileDescriptor_e5e70f8309ff9e66) }

var fileDescriptor_e5e70f8309ff9e66 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x91, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x31, 0x84, 0xd0, 0x5c, 0xb7, 0x43, 0x20, 0x2b, 0x42, 0x22, 0xca, 0x94, 0xc9, 0x11,
	0x65, 0x61, 0x44, 0xfc, 0x84, 0xd0, 0x89, 0xcd, 0x71, 0x4d, 0x15, 0xb5, 0x89, 0x83, 0xed, 0xc0,
	0xff, 0xe3, 0x97, 0xa1, 0xf8, 0x6a, 0x0a, 0x7b, 0xb7, 0xfb, 0xde, 0x3d, 0xe9, 0xe9, 0xde, 0xc1,
	0xcd, 0x97, 0xb1, 0x3b, 0x6d, 0xc7, 0xb6, 0xa6, 0x41, 0x8c, 0xd6, 0x78, 0x83, 0x29, 0x51, 0x7e,
	0xbf, 0x35, 0x66, 0xbb, 0xd7, 0x75, 0x50, 0xdb, 0xe9, 0xbd, 0xf6, 0x5d, 0xaf, 0x9d, 0x97, 0xfd,
	0x48, 0xc6, 0xd2, 0xc0, 0x72, 0x2d, 0xdd, 0xae, 0xd1, 0x1f, 0x93, 0x76, 0x1e, 0x11, 0x92, 0xc1,
	0x6c, 0x34, 0x67, 0x05, 0xab, 0xb2, 0x26, 0xcc, 0x78, 0x07, 0x99, 0x32, 0xc3, 0xa6, 0xf3, 0x9d,
	0x19, 0xf8, 0x79, 0x58, 0x1c, 0x05, 0xbc, 0x85, 0x54, 0xaa, 0xb0, 0xba, 0x08, 0xab, 0x03, 0xcd,
	0xfa, 0x28, 0xad, 0xec, 0x1d, 0x4f, 0x48, 0x27, 0x2a, 0xbf, 0x19, 0x00, 0x25, 0xba, 0x69, 0x7f,
	0xe2, 0x40, 0x3a, 0x3a, 0x06, 0x12, 0x21, 0x87, 0x2b, 0x37, 0x29, 0xa5, 0x9d, 0xe3, 0x97, 0x05,
	0xab, 0x16, 0x4d, 0x44, 0x7c, 0x82, 0xec, 0xb7, 0x0e, 0x9e, 0x16, 0xac, 0x5a, 0xae, 0x72, 0x41,
	0x85, 0x89, 0x58, 0x98, 0x58, 0x47, 0x47, 0x73, 0x34, 0xaf, 0x9e, 0xa9, 0xb5, 0x57, 0x6d, 0x3f,
	0x3b, 0xa5, 0xf1, 0x01, 0x92, 0x19, 0xf1, 0x5a, 0x1c, 0x9e, 0xf0, 0xa7, 0xd2, 0x1c, 0xff, 0x8b,
	0xf3, 0xd5, 0xe5, 0xd9, 0x0b, 0xbc, 0x2d, 0xe2, 0xe7, 0xda, 0x34, 0x84, 0x3d, 0xfe, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x6c, 0x62, 0x73, 0x0b, 0xcc, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TaskServiceClient interface {
	Task(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskResult, error)
}

type taskServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskServiceClient(cc *grpc.ClientConn) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) Task(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskResult, error) {
	out := new(TaskResult)
	err := c.cc.Invoke(ctx, "/worker.TaskService/Task", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
type TaskServiceServer interface {
	Task(context.Context, *TaskRequest) (*TaskResult, error)
}

// UnimplementedTaskServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (*UnimplementedTaskServiceServer) Task(ctx context.Context, req *TaskRequest) (*TaskResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Task not implemented")
}

func RegisterTaskServiceServer(s *grpc.Server, srv TaskServiceServer) {
	s.RegisterService(&_TaskService_serviceDesc, srv)
}

func _TaskService_Task_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).Task(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/worker.TaskService/Task",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).Task(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TaskService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "worker.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Task",
			Handler:    _TaskService_Task_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "workerpb/worker.proto",
}

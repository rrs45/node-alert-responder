// Code generated by protoc-gen-go. DO NOT EDIT.
// source: workerpb/worker.proto

package workerpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
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
	Source               string   `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`
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

func (m *TaskRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

type TaskAck struct {
	Condition            string   `protobuf:"bytes,1,opt,name=condition,proto3" json:"condition,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskAck) Reset()         { *m = TaskAck{} }
func (m *TaskAck) String() string { return proto.CompactTextString(m) }
func (*TaskAck) ProtoMessage()    {}
func (*TaskAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e70f8309ff9e66, []int{1}
}

func (m *TaskAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskAck.Unmarshal(m, b)
}
func (m *TaskAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskAck.Marshal(b, m, deterministic)
}
func (m *TaskAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskAck.Merge(m, src)
}
func (m *TaskAck) XXX_Size() int {
	return xxx_messageInfo_TaskAck.Size(m)
}
func (m *TaskAck) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskAck.DiscardUnknown(m)
}

var xxx_messageInfo_TaskAck proto.InternalMessageInfo

func (m *TaskAck) GetCondition() string {
	if m != nil {
		return m.Condition
	}
	return ""
}

type TaskStatus struct {
	Action               string   `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
	Worker               string   `protobuf:"bytes,2,opt,name=worker,proto3" json:"worker,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskStatus) Reset()         { *m = TaskStatus{} }
func (m *TaskStatus) String() string { return proto.CompactTextString(m) }
func (*TaskStatus) ProtoMessage()    {}
func (*TaskStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e70f8309ff9e66, []int{2}
}

func (m *TaskStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskStatus.Unmarshal(m, b)
}
func (m *TaskStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskStatus.Marshal(b, m, deterministic)
}
func (m *TaskStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskStatus.Merge(m, src)
}
func (m *TaskStatus) XXX_Size() int {
	return xxx_messageInfo_TaskStatus.Size(m)
}
func (m *TaskStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskStatus.DiscardUnknown(m)
}

var xxx_messageInfo_TaskStatus proto.InternalMessageInfo

func (m *TaskStatus) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *TaskStatus) GetWorker() string {
	if m != nil {
		return m.Worker
	}
	return ""
}

type AllTasks struct {
	Items                map[string]*TaskStatus `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *AllTasks) Reset()         { *m = AllTasks{} }
func (m *AllTasks) String() string { return proto.CompactTextString(m) }
func (*AllTasks) ProtoMessage()    {}
func (*AllTasks) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e70f8309ff9e66, []int{3}
}

func (m *AllTasks) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AllTasks.Unmarshal(m, b)
}
func (m *AllTasks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AllTasks.Marshal(b, m, deterministic)
}
func (m *AllTasks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllTasks.Merge(m, src)
}
func (m *AllTasks) XXX_Size() int {
	return xxx_messageInfo_AllTasks.Size(m)
}
func (m *AllTasks) XXX_DiscardUnknown() {
	xxx_messageInfo_AllTasks.DiscardUnknown(m)
}

var xxx_messageInfo_AllTasks proto.InternalMessageInfo

func (m *AllTasks) GetItems() map[string]*TaskStatus {
	if m != nil {
		return m.Items
	}
	return nil
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
	return fileDescriptor_e5e70f8309ff9e66, []int{4}
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
	proto.RegisterType((*TaskAck)(nil), "worker.TaskAck")
	proto.RegisterType((*TaskStatus)(nil), "worker.TaskStatus")
	proto.RegisterType((*AllTasks)(nil), "worker.AllTasks")
	proto.RegisterMapType((map[string]*TaskStatus)(nil), "worker.AllTasks.ItemsEntry")
	proto.RegisterType((*TaskResult)(nil), "worker.TaskResult")
}

func init() { proto.RegisterFile("workerpb/worker.proto", fileDescriptor_e5e70f8309ff9e66) }

var fileDescriptor_e5e70f8309ff9e66 = []byte{
	// 436 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xd5, 0x36, 0x1f, 0x4d, 0x26, 0x20, 0xca, 0x20, 0x90, 0xe5, 0x22, 0x51, 0xf9, 0x42, 0x0e,
	0xc8, 0x11, 0xee, 0xa5, 0x7c, 0x5c, 0x82, 0x54, 0xa1, 0x4a, 0x88, 0x83, 0x5b, 0x2e, 0xdc, 0x9c,
	0xcd, 0x50, 0x59, 0xb6, 0xb3, 0xc6, 0xbb, 0x0e, 0xca, 0x2f, 0xe0, 0xc0, 0xbf, 0xe2, 0x97, 0xa1,
	0xdd, 0xb1, 0xb1, 0xeb, 0x5e, 0x7b, 0xdb, 0xf7, 0xe6, 0x69, 0xde, 0x9b, 0xf1, 0x18, 0x9e, 0xff,
	0x52, 0x55, 0x46, 0x55, 0xb9, 0x59, 0xf1, 0x23, 0x2c, 0x2b, 0x65, 0x14, 0x4e, 0x19, 0xf9, 0xaf,
	0x6e, 0x95, 0xba, 0xcd, 0x69, 0xe5, 0xd8, 0x4d, 0xfd, 0x63, 0x65, 0xd2, 0x82, 0xb4, 0x49, 0x8a,
	0x92, 0x85, 0xfe, 0xe9, 0x50, 0x40, 0x45, 0x69, 0x0e, 0x5c, 0x0c, 0x7e, 0x0b, 0x58, 0xdc, 0x24,
	0x3a, 0x8b, 0xe9, 0x67, 0x4d, 0xda, 0x20, 0xc2, 0x78, 0xa7, 0xb6, 0xe4, 0x89, 0x33, 0xb1, 0x9c,
	0xc7, 0xee, 0x8d, 0x2f, 0x61, 0x2e, 0xd5, 0x6e, 0x9b, 0x9a, 0x54, 0xed, 0xbc, 0x23, 0x57, 0xe8,
	0x08, 0x7c, 0x01, 0xd3, 0x44, 0xba, 0xd2, 0xc8, 0x95, 0x1a, 0x64, 0xf9, 0x32, 0xa9, 0x92, 0x42,
	0x7b, 0x63, 0xe6, 0x19, 0x59, 0x5e, 0xab, 0xba, 0x92, 0xe4, 0x4d, 0x98, 0x67, 0x14, 0xbc, 0x86,
	0x63, 0x1b, 0x64, 0x2d, 0xb3, 0xbb, 0x86, 0x62, 0x60, 0x18, 0x7c, 0x04, 0xb0, 0xc2, 0x6b, 0x93,
	0x98, 0x5a, 0xf7, 0xec, 0xc5, 0xd0, 0x9e, 0x17, 0xd4, 0x24, 0x6e, 0x50, 0xf0, 0x47, 0xc0, 0x6c,
	0x9d, 0xe7, 0xb6, 0x83, 0xc6, 0xb7, 0x30, 0x49, 0x0d, 0x15, 0xda, 0x13, 0x67, 0xa3, 0xe5, 0x22,
	0x3a, 0x0d, 0x9b, 0x0d, 0xb7, 0x82, 0xf0, 0xca, 0x56, 0x2f, 0x77, 0xa6, 0x3a, 0xc4, 0xac, 0xf4,
	0xbf, 0x00, 0x74, 0x24, 0x9e, 0xc0, 0x28, 0xa3, 0x43, 0x63, 0x6d, 0x9f, 0xb8, 0x84, 0xc9, 0x3e,
	0xc9, 0x6b, 0x72, 0xb6, 0x8b, 0x08, 0xdb, 0x96, 0x5d, 0xe4, 0x98, 0x05, 0xef, 0x8f, 0x2e, 0x44,
	0xf0, 0x57, 0xf0, 0x30, 0x31, 0xe9, 0x3a, 0x7f, 0xe0, 0xed, 0x37, 0xe3, 0x8f, 0xfb, 0xe3, 0xa3,
	0x07, 0xc7, 0xba, 0x96, 0x92, 0xb4, 0x76, 0xeb, 0x9f, 0xc5, 0x2d, 0xc4, 0x0b, 0x98, 0xff, 0xbf,
	0x1c, 0x6f, 0xea, 0xc2, 0xfb, 0x21, 0x9f, 0x4e, 0xd8, 0x9e, 0x4e, 0x78, 0xd3, 0x2a, 0xe2, 0x4e,
	0x1c, 0x7d, 0xe0, 0x13, 0xba, 0xa6, 0x6a, 0x9f, 0x4a, 0xc2, 0x37, 0x30, 0xb6, 0x10, 0x9f, 0xf5,
	0x47, 0x6f, 0xee, 0xcb, 0x7f, 0xd2, 0x27, 0xd7, 0x32, 0x8b, 0xbe, 0xc2, 0xd3, 0x6e, 0x35, 0x6d,
	0x8b, 0x77, 0xf0, 0xf8, 0x33, 0x99, 0xfe, 0x57, 0xbe, 0x97, 0xe4, 0xd2, 0x1e, 0xb1, 0x7f, 0x32,
	0xfc, 0x62, 0xd1, 0x15, 0x20, 0xfb, 0x49, 0x4a, 0xf7, 0xd4, 0x36, 0x3c, 0x87, 0x47, 0xbc, 0xe2,
	0x6f, 0xe5, 0x36, 0x31, 0x84, 0x78, 0x37, 0x9b, 0xad, 0xdc, 0x8b, 0xf6, 0x09, 0xbe, 0xcf, 0xda,
	0x5f, 0x6f, 0x33, 0x75, 0xc6, 0xe7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x27, 0x39, 0x5a, 0x74,
	0x8d, 0x03, 0x00, 0x00,
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
	Task(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskAck, error)
}

type taskServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskServiceClient(cc *grpc.ClientConn) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) Task(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskAck, error) {
	out := new(TaskAck)
	err := c.cc.Invoke(ctx, "/worker.TaskService/Task", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
type TaskServiceServer interface {
	Task(context.Context, *TaskRequest) (*TaskAck, error)
}

// UnimplementedTaskServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (*UnimplementedTaskServiceServer) Task(ctx context.Context, req *TaskRequest) (*TaskAck, error) {
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

// TaskStatusServiceClient is the client API for TaskStatusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TaskStatusServiceClient interface {
	GetTaskStatus(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*AllTasks, error)
}

type taskStatusServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskStatusServiceClient(cc *grpc.ClientConn) TaskStatusServiceClient {
	return &taskStatusServiceClient{cc}
}

func (c *taskStatusServiceClient) GetTaskStatus(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*AllTasks, error) {
	out := new(AllTasks)
	err := c.cc.Invoke(ctx, "/worker.TaskStatusService/GetTaskStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskStatusServiceServer is the server API for TaskStatusService service.
type TaskStatusServiceServer interface {
	GetTaskStatus(context.Context, *empty.Empty) (*AllTasks, error)
}

// UnimplementedTaskStatusServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTaskStatusServiceServer struct {
}

func (*UnimplementedTaskStatusServiceServer) GetTaskStatus(ctx context.Context, req *empty.Empty) (*AllTasks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStatus not implemented")
}

func RegisterTaskStatusServiceServer(s *grpc.Server, srv TaskStatusServiceServer) {
	s.RegisterService(&_TaskStatusService_serviceDesc, srv)
}

func _TaskStatusService_GetTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskStatusServiceServer).GetTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/worker.TaskStatusService/GetTaskStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskStatusServiceServer).GetTaskStatus(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _TaskStatusService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "worker.TaskStatusService",
	HandlerType: (*TaskStatusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTaskStatus",
			Handler:    _TaskStatusService_GetTaskStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "workerpb/worker.proto",
}

// TaskReceiveServiceClient is the client API for TaskReceiveService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TaskReceiveServiceClient interface {
	ResultUpdate(ctx context.Context, in *TaskResult, opts ...grpc.CallOption) (*TaskAck, error)
}

type taskReceiveServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskReceiveServiceClient(cc *grpc.ClientConn) TaskReceiveServiceClient {
	return &taskReceiveServiceClient{cc}
}

func (c *taskReceiveServiceClient) ResultUpdate(ctx context.Context, in *TaskResult, opts ...grpc.CallOption) (*TaskAck, error) {
	out := new(TaskAck)
	err := c.cc.Invoke(ctx, "/worker.TaskReceiveService/ResultUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskReceiveServiceServer is the server API for TaskReceiveService service.
type TaskReceiveServiceServer interface {
	ResultUpdate(context.Context, *TaskResult) (*TaskAck, error)
}

// UnimplementedTaskReceiveServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTaskReceiveServiceServer struct {
}

func (*UnimplementedTaskReceiveServiceServer) ResultUpdate(ctx context.Context, req *TaskResult) (*TaskAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResultUpdate not implemented")
}

func RegisterTaskReceiveServiceServer(s *grpc.Server, srv TaskReceiveServiceServer) {
	s.RegisterService(&_TaskReceiveService_serviceDesc, srv)
}

func _TaskReceiveService_ResultUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskResult)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskReceiveServiceServer).ResultUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/worker.TaskReceiveService/ResultUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskReceiveServiceServer).ResultUpdate(ctx, req.(*TaskResult))
	}
	return interceptor(ctx, in, info, handler)
}

var _TaskReceiveService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "worker.TaskReceiveService",
	HandlerType: (*TaskReceiveServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResultUpdate",
			Handler:    _TaskReceiveService_ResultUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "workerpb/worker.proto",
}

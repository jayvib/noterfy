package protoutil

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"noterfy/note"
	pb "noterfy/note/proto"
)

var errUnexpected = errors.New("unexpected write count")

// WriteProtoMessage marshals the m protocol buffer then writes to
// w writer. It prepend a 4-byte size prior to protobuf binary.
func WriteProtoMessage(w io.Writer, message proto.Message) error {

	msgBytes, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	msgLen := len(msgBytes)

	// Create a 4-byte binary that contains
	// the length of the message. Use little endian.
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(msgLen))

	n, err := w.Write(buf)
	if err != nil {
		return err
	}

	if n != 4 {
		return errUnexpected
	}

	// Write to protocol buffer binary
	n, err = w.Write(msgBytes)
	if err != nil {
		return err
	}

	if n != msgLen {
		return errUnexpected
	}

	return nil
}

// WriteAllProtoMessages is like a WriteProtoMessage but it takes an arbitrary
// messages to write to w writer.
func WriteAllProtoMessages(w io.Writer, messages ...proto.Message) error {
	for _, message := range messages {
		err := WriteProtoMessage(w, message)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadProtoMessage reads protobuf encoded content from r.
// It un-marshals the content into a note protobuf message then
// returns the note. If there's an error it could be an io.EOF error.
func ReadProtoMessage(r io.Reader) (*note.Note, error) {
	msgLen := make([]byte, 4)
	_, err := io.ReadFull(r, msgLen)
	if err != nil {
		return nil, err
	}

	size := binary.LittleEndian.Uint32(msgLen)
	gotSize := int(size)

	msg := make([]byte, gotSize)
	_, err = io.ReadFull(r, msg)
	if err != nil {
		return nil, err
	}

	var got pb.Note
	err = proto.Unmarshal(msg, &got)
	if err != nil {
		return nil, err
	}

	return ProtoToNote(&got)
}

// ReadAllProtoMessages reads all the proto messages from r until
// the reader return an io.EOF.
func ReadAllProtoMessages(r io.Reader) ([]*note.Note, error) {

	var notes []*note.Note
	for {
		n, err := ReadProtoMessage(r)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, nil
}

// ProtoToNote converts the note protocol buffer message
// to note.Note. If there's any error, it will be related
// to UUID byte parsing.
func ProtoToNote(p *pb.Note) (*note.Note, error) {
	id, err := uuid.ParseBytes(p.Id)
	if err != nil {
		return nil, err
	}
	n := new(note.Note)
	n.SetID(id).
		SetTitle(p.Title).
		SetContent(p.Content).
		SetCreatedTime(p.CreatedTime.AsTime()).
		SetUpdatedTime(p.UpdatedTime.AsTime()).
		SetIsFavorite(p.IsFavorite)
	return n, nil
}

// NoteToProto converts the note to protocol buffer message.
func NoteToProto(n *note.Note) *pb.Note {
	return &pb.Note{
		Id:          []byte(n.ID.String()),
		Title:       n.GetTitle(),
		Content:     n.GetContent(),
		CreatedTime: timestamppb.New(n.GetCreatedTime()),
		UpdatedTime: timestamppb.New(n.GetUpdatedTime()),
		IsFavorite:  n.GetIsFavorite(),
	}
}

// ConvertNotesToProtos convert the array of notes into a
// note protocol buffer message.
func ConvertNotesToProtos(notes []*note.Note) (pbs []*pb.Note) {
	for _, n := range notes {
		pbs = append(pbs, NoteToProto(n))
	}
	return
}

// ConvertToProtoMessage takes an array of protobuf note message and
// converts it to proto.Message interface.
func ConvertToProtoMessage(pbs []*pb.Note) []proto.Message {
	var messages []proto.Message
	for _, p := range pbs {
		messages = append(messages, p)
	}
	return messages
}

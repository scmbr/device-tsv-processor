package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scmbr/device-tsv-processor/internal/domain"
	"github.com/scmbr/device-tsv-processor/internal/errs"
	"github.com/scmbr/device-tsv-processor/internal/usecase/device_message"
)

type DeviceMessagesHandler struct {
	get_uc device_message.GetDeviceMessages
}

func NewDeviceMessageHandler(get_uc device_message.GetDeviceMessages) *DeviceMessagesHandler {
	return &DeviceMessagesHandler{get_uc: get_uc}
}

type DeviceMessageResponse struct {
	ID        int64     `json:"id"`
	DeviceID  int64     `json:"device_id"`
	InvID     string    `json:"inv_id"`
	MsgID     string    `json:"msg_id"`
	Text      string    `json:"text"`
	Context   string    `json:"context"`
	Class     string    `json:"class"`
	Level     int       `json:"level"`
	Area      string    `json:"area"`
	Addr      string    `json:"addr"`
	Block     string    `json:"block"`
	Type      string    `json:"type"`
	Bit       int       `json:"bit"`
	InvertBit bool      `json:"invert_bit"`
	CreatedAt time.Time `json:"created_at"`
}
type ListMessagesByGUIDResponse struct {
	Total    int                     `json:"total"`
	Messages []DeviceMessageResponse `json:"messages"`
}

func toMessageResponse(m *domain.DeviceMessage) DeviceMessageResponse {
	return DeviceMessageResponse{
		ID:        m.ID,
		DeviceID:  m.DeviceID,
		InvID:     m.InvID,
		MsgID:     m.MsgID,
		Text:      m.Text,
		Context:   m.Context,
		Class:     m.Class,
		Level:     m.Level,
		Area:      m.Area,
		Addr:      m.Addr,
		Block:     m.Block,
		Type:      m.Type,
		Bit:       m.Bit,
		InvertBit: m.InvertBit,
		CreatedAt: m.CreatedAt,
	}
}
func (h *DeviceMessagesHandler) ListByGUID(ctx *gin.Context) {
	const op = "device_message.http.list_by_guid"
	guid := ctx.Param("guid")
	if guid == "" {
		ctx.Error(errs.E(errs.KindInvalid, "INVALID_GUID", op, "invalid guid", map[string]string{"guid": "is required"}, nil))
		return
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	if err != nil {
		ctx.Error(errs.E(errs.KindInvalid, "INVALID_LIMIT", op, "invalid limit", map[string]string{"limit": "must have type integer"}, nil))
		return
	}
	if limit > 100 || limit <= 0 {
		ctx.Error(errs.E(errs.KindInvalid, "INVALID_LIMIT", op, "invalid limit", map[string]string{"limit": "is not valid"}, nil))
		return
	}
	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		ctx.Error(errs.E(errs.KindInvalid, "INVALID_OFFSET", op, "invalid offset", map[string]string{"offset": "must have type integer"}, nil))
		return
	}
	if offset < 0 {
		ctx.Error(errs.E(errs.KindInvalid, "INVALID_OFFSET", op, "invalid offset", map[string]string{"offset": "is not valid"}, nil))
		return
	}
	messages, total, err := h.get_uc.Execute(ctx.Request.Context(), device_message.GetDeviceMessagesInput{
		Limit:  limit,
		Offset: offset,
		GUID:   guid,
	})
	if err != nil {
		ctx.Error(err)
		return
	}
	messagesResponse := make([]DeviceMessageResponse, 0, len(messages))
	for _, m := range messages {
		messagesResponse = append(messagesResponse, toMessageResponse(m))
	}
	out := ListMessagesByGUIDResponse{
		Total:    total,
		Messages: messagesResponse,
	}
	ctx.JSON(http.StatusOK, out)
}

package api

import (
	"mobile-connect/internal/stream"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) StreamVideo(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	for {

		frame, err := h.streamService.Capture(id)

		if err != nil {
			break
		}

		err = conn.WriteMessage(
			websocket.BinaryMessage,
			frame,
		)

		if err != nil {
			break
		}
	}
}

func (h *Handler) StreamControl(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	for {

		var event stream.InputEvent

		err := conn.ReadJSON(&event)

		if err != nil {
			break
		}

		switch event.Type {

		case "tap":

			_ = h.streamService.Tap(
				id,
				int(event.X),
				int(event.Y),
			)

		case "swipe":

			_ = h.streamService.Swipe(
				id,
				int(event.X),
				int(event.Y),
				int(event.X2),
				int(event.Y2),
				event.Duration,
			)

		case "longpress":

			_ = h.streamService.LongPress(
				id,
				int(event.X),
				int(event.Y),
				event.Duration,
			)

		case "text":

			_ = h.streamService.InputText(
				id,
				event.Text,
			)
		}
	}
}

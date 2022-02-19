package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch r.Method {
	case http.MethodGet:
		size := r.URL.Query().Get("size")
		var err error
		size64 := int64(5)
		if size != "" {
			size64, err = strconv.ParseInt(size, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		prevId := r.URL.Query().Get("prev_id")
		prevId64 := int64(0)
		if prevId != "" {
			prevId64, err = strconv.ParseInt(prevId, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
		}
		request := &model.ReadTODORequest{Size: size64, PrevID: prevId64}

		response, err := h.Read(ctx, request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)

		if err := encoder.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var request model.CreateTODORequest
		err := decoder.Decode(&request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Subject == "" {
			log.Println("Subject not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Create(ctx, &request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)

		err = encoder.Encode(res)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		var request model.UpdateTODORequest
		err := decoder.Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		if request.ID == 0 {
			log.Println("ID not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Subject == "" {
			log.Println("Subject not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := h.Update(ctx, &request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")

		err = encoder.Encode(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

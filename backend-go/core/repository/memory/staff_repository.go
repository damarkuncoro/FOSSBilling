package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockStaffRepository struct {
	mu        sync.RWMutex
	staffs    map[int64]*domain.Staff
	groups    map[int64]*domain.AdminGroup
	auditLogs []*domain.AuditLog
	nextID    int64
	nextGrpID int64
	nextLogID int64
}

func NewMockStaffRepository() *MockStaffRepository {
	return &MockStaffRepository{
		staffs:    make(map[int64]*domain.Staff),
		groups:    make(map[int64]*domain.AdminGroup),
		nextID:    1,
		nextGrpID: 1,
		nextLogID: 1,
	}
}

func (r *MockStaffRepository) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.staffs[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *s
	return &cp, nil
}

func (r *MockStaffRepository) GetByEmail(ctx context.Context, email string) (*domain.Staff, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, s := range r.staffs {
		if strings.EqualFold(s.Email, email) {
			cp := *s
			return &cp, nil
		}
	}
	return nil, appErrors.ErrNotFound
}

func (r *MockStaffRepository) List(ctx context.Context, limit, offset int) ([]*domain.Staff, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.Staff
	for _, s := range r.staffs {
		cp := *s
		list = append(list, &cp)
	}

	total := len(list)
	if offset >= total {
		return []*domain.Staff{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return list[offset:end], total, nil
}

func (r *MockStaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	staff.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	staff.CreatedAt = now
	staff.UpdatedAt = now
	if staff.Status == "" {
		staff.Status = "active"
	}

	cp := *staff
	r.staffs[staff.ID] = &cp
	return nil
}

func (r *MockStaffRepository) Update(ctx context.Context, staff *domain.Staff) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.staffs[staff.ID]; !ok {
		return appErrors.ErrNotFound
	}
	staff.UpdatedAt = time.Now().UTC()
	cp := *staff
	r.staffs[staff.ID] = &cp
	return nil
}

func (r *MockStaffRepository) GetGroupByID(ctx context.Context, groupID int64) (*domain.AdminGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	g, ok := r.groups[groupID]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *g
	return &cp, nil
}

func (r *MockStaffRepository) CreateGroup(ctx context.Context, group *domain.AdminGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	group.ID = r.nextGrpID
	r.nextGrpID++
	now := time.Now().UTC()
	group.CreatedAt = now
	group.UpdatedAt = now

	cp := *group
	r.groups[group.ID] = &cp
	return nil
}

func (r *MockStaffRepository) AddAuditLog(ctx context.Context, log *domain.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.ID = r.nextLogID
	r.nextLogID++
	log.CreatedAt = time.Now().UTC()

	cp := *log
	r.auditLogs = append(r.auditLogs, &cp)
	return nil
}

func (r *MockStaffRepository) ListAuditLogs(ctx context.Context, limit, offset int) ([]*domain.AuditLog, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.auditLogs)
	if offset >= total {
		return []*domain.AuditLog{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	// Try to populate staff names from the mock state
	res := r.auditLogs[offset:end]
	for _, l := range res {
		if l.StaffID != nil {
			if s, ok := r.staffs[*l.StaffID]; ok {
				l.StaffName = &s.Name
			}
		}
	}

	return res, total, nil
}

func (r *MockStaffRepository) GetAuditLogs() []*domain.AuditLog {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.auditLogs
}

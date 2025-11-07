package access

import "testing"

func TestCanAccess(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		resource string
		expected bool
	}{
		{
			name:     "Admin có quyền truy cập mọi tài nguyên",
			role:     "admin",
			resource: "admin_panel",
			expected: true,
		},
		{
			name:     "User có quyền truy cập tài nguyên bình thường",
			role:     "user",
			resource: "dashboard",
			expected: true,
		},
		{
			name:     "User không được truy cập admin_panel",
			role:     "user",
			resource: "admin_panel",
			expected: false,
		},
		{
			name:     "Khách không có quyền truy cập",
			role:     "guest",
			resource: "dashboard",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanAccess(tt.role, tt.resource)
			if got != tt.expected {
				t.Errorf("CanAccess(%q, %q) = %v, want %v",
					tt.role, tt.resource, got, tt.expected)
			}
		})
	}
}

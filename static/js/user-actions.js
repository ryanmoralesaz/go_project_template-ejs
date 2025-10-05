// User management actions (delete, etc.)

async function deleteUser(userId) {
  console.log('Deleting user:', userId);

  if (!confirm('Are you sure you want to delete this user? This action cannot be undone.')) {
    return;
  }

  // Find the button that was clicked to provide visual feedback
  const deleteButton = event.target.closest('button');
  const originalText = deleteButton.innerHTML;

  try {
    // Show loading state
    deleteButton.disabled = true;
    deleteButton.innerHTML = '<i class="bi bi-spinner-border spinner-border-sm me-1"></i>Deleting...';

    const response = await fetch(`/users/${userId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
    });

    const data = await response.json();
    console.log('Delete response:', data);

    if (data.success) {
      // Show success message
      deleteButton.innerHTML = '<i class="bi bi-check-circle me-1"></i>Deleted!';
      deleteButton.classList.remove('btn-danger');
      deleteButton.classList.add('btn-success');

      // Remove the user card from the DOM with animation
      const userCard = deleteButton.closest('.col-md-6, .col-lg-4');
      if (userCard) {
        userCard.style.transition = 'opacity 0.3s ease-out, transform 0.3s ease-out';
        userCard.style.opacity = '0';
        userCard.style.transform = 'scale(0.9)';

        setTimeout(() => {
          userCard.remove();

          // Check if there are no more users and show appropriate message
          const remainingCards = document.querySelectorAll('.col-md-6, .col-lg-4');
          if (remainingCards.length === 0) {
            location.reload(); // Reload to show "no users" message
          }
        }, 300);
      }

    } else {
      throw new Error(data.message || 'Failed to delete user');
    }
  } catch (err) {
    console.error('Error deleting user:', err);

    // Show error message
    alert(`Failed to delete user: ${err.message}`);

    // Reset button state
    deleteButton.disabled = false;
    deleteButton.innerHTML = originalText;
    deleteButton.classList.remove('btn-success');
    deleteButton.classList.add('btn-danger');
  }
}

// Utility function to show toast notifications (if you want to add this later)
function showToast(message, type = 'info') {
  // Create toast element
  const toast = document.createElement('div');
  toast.className = `alert alert-${type} position-fixed`;
  toast.style.cssText = `
    top: 20px;
    right: 20px;
    z-index: 9999;
    min-width: 300px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  `;
  toast.innerHTML = `
    <div class="d-flex align-items-center">
      <i class="bi bi-${type === 'success' ? 'check-circle' : type === 'danger' ? 'exclamation-triangle' : 'info-circle'} me-2"></i>
      ${message}
      <button type="button" class="btn-close ms-auto" onclick="this.parentElement.parentElement.remove()"></button>
    </div>
  `;

  // Add to page
  document.body.appendChild(toast);

  // Auto-remove after 5 seconds
  setTimeout(() => {
    if (toast.parentElement) {
      toast.remove();
    }
  }, 5000);
}
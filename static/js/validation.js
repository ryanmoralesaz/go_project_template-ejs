document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('registrationForm');
  if (!form) return; // Exit if form doesn't exist on this page

  const validationState = {};
  const inputs = form.querySelectorAll('input[pattern]');

  function setErrorMessage(fieldName, value, errorElement) {
    let errorMessage = '';

    switch (fieldName) {
      case 'firstName':
        if (value.length === 0) {
          errorMessage = 'First name is required';
        } else if (value.length > 50) {
          errorMessage = 'First name must be less than 50 characters';
        } else if (value.length < 2) {
          errorMessage = 'First name must be 2 or more characters';
        } else {
          errorMessage = 'First name must only include letters';
        }
        break;

      case 'lastName':
        if (value.length === 0) {
          errorMessage = 'Last name is required';
        } else if (value.length > 50) {
          errorMessage = 'Last name must be less than 50 characters';
        } else if (value.length < 2) {
          errorMessage = 'Last name must be 2 or more characters';
        } else {
          errorMessage = 'Last name must only include letters';
        }
        break;

      case 'email':
        if (value.length === 0) {
          errorMessage = 'Email address is required';
        } else {
          errorMessage = 'Please enter a valid email address';
        }
        break;

      case 'phone':
        if (value.length === 0) {
          errorMessage = 'Phone number is required';
        } else {
          errorMessage = 'Please enter a valid phone number (format: 123-456-7890)';
        }
        break;

      default:
        errorMessage = `Invalid ${fieldName}`;
    }

    errorElement.textContent = errorMessage;
  }

  function validateInput(input) {
    const fieldName = input.name;
    const value = input.value.trim();
    const pattern = new RegExp(input.pattern);
    const isValid = pattern.test(value);
    const errorElement = document.getElementById(`${fieldName}Error`);
    const validationIcon = input.parentElement.querySelector('.validation-icon');

    console.log(`${fieldName} validation:`, isValid);

    if (isValid) {
      input.classList.add('is-valid');
      input.classList.remove('is-invalid');
      if (validationIcon) {
        validationIcon.style.display = 'block';
      }
      errorElement.textContent = '';
    } else {
      input.classList.remove('is-valid');
      input.classList.add('is-invalid');
      if (validationIcon) {
        validationIcon.style.display = 'none';
      }
      setErrorMessage(fieldName, value, errorElement);
    }

    validationState[fieldName] = isValid;
    console.log('Validation state:', validationState);
  }

  // Initialize validation for all inputs
  inputs.forEach((input) => {
    const fieldName = input.name;
    validationState[fieldName] = false;

    // Validate on input (real-time)
    input.addEventListener('input', function () {
      validateInput(this);
    });

    // Validate on blur (when user leaves field)
    input.addEventListener('blur', function () {
      validateInput(this);
    });
  });

  // Handle form submission
  form.addEventListener('submit', async function (e) {
    e.preventDefault();

    const formData = new FormData(this);
    const requiredFields = document.querySelectorAll('input[required]');

    // Validate all fields before submission
    inputs.forEach(validateInput);

    // Check if form is valid
    const isFormValid = [...requiredFields].every((input) => {
      const value = formData.get(input.name);
      return value && validationState[input.name];
    });

    if (!isFormValid) {
      console.log('Form invalid - missing required fields or validation failed');

      // Focus on first invalid field
      const firstInvalidField = form.querySelector('.is-invalid');
      if (firstInvalidField) {
        firstInvalidField.focus();
      }

      return;
    }

    // Collect all form data
    const submitData = {};
    for (const [key, value] of formData.entries()) {
      if (value.trim()) {
        submitData[key] = value.trim();
      }
    }

    console.log('Submitting:', submitData);

    // Disable submit button during submission
    const submitBtn = form.querySelector('#submitBtn');
    const originalText = submitBtn.innerHTML;
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<i class="bi bi-spinner-border spinner-border-sm me-1"></i>Submitting...';

    try {
      const response = await fetch('/form/new', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify(submitData),
      });

      const data = await response.json();
      console.log('Response from server:', data);

      if (data.success && data.redirectTo) {
        // Show success message briefly before redirect
        submitBtn.innerHTML = '<i class="bi bi-check-circle me-1"></i>Success!';
        submitBtn.classList.remove('btn-primary');
        submitBtn.classList.add('btn-success');

        setTimeout(() => {
          window.location.href = data.redirectTo;
        }, 1000);
      } else {
        throw new Error(data.message || 'Unknown error occurred');
      }
    } catch (err) {
      console.error('Submission error:', err);

      // Show error message
      alert(`Error: ${err.message}`);

      // Reset button
      submitBtn.disabled = false;
      submitBtn.innerHTML = originalText;
    }
  });
});
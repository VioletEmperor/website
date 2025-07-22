// Firebase functions will be loaded from window object
let auth, onAuthStateChanged, signOut;
// Auth state management
let currentUser = null;
let authToken = null;

// DOM elements
let adminEmailElement, settingsEmailElement, logoutButton;
let navTabs, tabContents, fileUpload, fileInput, fileInfo;

document.addEventListener('DOMContentLoaded', function() {
  // Initialize DOM elements
  adminEmailElement = document.getElementById('admin-email');
  settingsEmailElement = document.getElementById('settings-email');
  logoutButton = document.querySelector('.btn-logout');
  navTabs = document.querySelectorAll('.nav-tab');
  tabContents = document.querySelectorAll('.tab-content');
  fileUpload = document.getElementById('file-upload');
  fileInput = document.getElementById('file-input');
  fileInfo = document.getElementById('file-info');

  // Wait for Firebase to be ready
  if (window.firebaseReady) {
    initializeFirebase();
  } else {
    window.addEventListener('firebaseReady', initializeFirebase);
  }
});

function initializeFirebase() {
  // Get Firebase functions from window object
  auth = window.firebaseAuth;
  onAuthStateChanged = window.onAuthStateChanged;
  signOut = window.signOut;

  // Initialize authentication
  initializeAuth();
  
  // Set up event listeners
  setupEventListeners();
}

function initializeAuth() {
  // Check for stored auth token
  const storedToken = localStorage.getItem('adminToken');
  if (storedToken) {
    authToken = storedToken;
  }

  // Listen for auth state changes
  onAuthStateChanged(auth, async (user) => {
    if (user) {
      currentUser = user;
      authToken = await user.getIdToken();
      localStorage.setItem('adminToken', authToken);
      
      // Update UI with user email
      updateUserEmail(user.email);
      
      // Verify admin access
      const isAdmin = await verifyAdminAccess(user.email);
      if (!isAdmin) {
        alert('You are not authorized to access the admin area.');
        await signOutUser();
      }
    } else {
      currentUser = null;
      authToken = null;
      localStorage.removeItem('adminToken');
      
      // Redirect to login
      window.location.href = '/admin/login';
    }
  });
}

function setupEventListeners() {
  // Tab switching
  navTabs.forEach(tab => {
    tab.addEventListener('click', () => {
      const targetTab = tab.getAttribute('data-tab');
      switchTab(targetTab);
    });
  });

  // File upload functionality
  if (fileUpload && fileInput) {
    setupFileUpload();
  }

  // Delete confirmation
  const deleteButtons = document.querySelectorAll('.btn-delete');
  deleteButtons.forEach(button => {
    button.addEventListener('click', handleDeletePost);
  });

  // Logout functionality
  if (logoutButton) {
    logoutButton.addEventListener('click', handleLogout);
  }

  // Form submissions
  const uploadForm = document.querySelector('.upload-form');
  if (uploadForm) {
    uploadForm.addEventListener('submit', handleFormSubmit);
  }
}

function switchTab(targetTab) {
  // Remove active class from all tabs and contents
  navTabs.forEach(t => t.classList.remove('active'));
  tabContents.forEach(content => content.classList.remove('active'));

  // Add active class to clicked tab and corresponding content
  document.querySelector(`[data-tab="${targetTab}"]`).classList.add('active');
  document.getElementById(targetTab).classList.add('active');
}

function setupFileUpload() {
  // Click to upload
  fileUpload.addEventListener('click', () => {
    fileInput.click();
  });

  // Drag and drop
  fileUpload.addEventListener('dragover', (e) => {
    e.preventDefault();
    fileUpload.classList.add('dragover');
  });

  fileUpload.addEventListener('dragleave', () => {
    fileUpload.classList.remove('dragover');
  });

  fileUpload.addEventListener('drop', (e) => {
    e.preventDefault();
    fileUpload.classList.remove('dragover');
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      fileInput.files = files;
      showFileInfo(files[0]);
    }
  });

  // File selection
  fileInput.addEventListener('change', () => {
    if (fileInput.files.length > 0) {
      showFileInfo(fileInput.files[0]);
    }
  });
}

function showFileInfo(file) {
  if (fileInfo) {
    fileInfo.textContent = `Selected: ${file.name} (${(file.size / 1024).toFixed(1)} KB)`;
    fileInfo.style.display = 'block';
  }
}

function handleDeletePost(e) {
  e.preventDefault();
  const postTitle = e.target.getAttribute('data-title') || 'this post';
  
  if (confirm(`Are you sure you want to delete "${postTitle}"? This action cannot be undone.`)) {
    // TODO: Implement delete functionality with API call
    window.location.href = e.target.href;
  }
}

async function handleLogout(e) {
  e.preventDefault();
  
  if (confirm('Are you sure you want to logout?')) {
    await signOutUser();
  }
}

async function handleFormSubmit(e) {
  // Add auth token to form data
  if (authToken) {
    // Create hidden input for auth token
    const tokenInput = document.createElement('input');
    tokenInput.type = 'hidden';
    tokenInput.name = 'authToken';
    tokenInput.value = authToken;
    e.target.appendChild(tokenInput);
  }
}

async function verifyAdminAccess(email) {
  try {
    if (!currentUser) return false;
    
    const idToken = await currentUser.getIdToken();
    
    const response = await fetch('/admin/verify', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${idToken}`
      },
      body: JSON.stringify({ email })
    });
    
    return response.ok;
  } catch (error) {
    console.error('Admin verification error:', error);
    return false;
  }
}

async function signOutUser() {
  try {
    await signOut(auth);
    localStorage.removeItem('adminToken');
    // Clear cookie
    document.cookie = 'adminToken=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
    window.location.href = '/admin/login';
  } catch (error) {
    console.error('Sign out error:', error);
  }
}

function updateUserEmail(email) {
  if (adminEmailElement) {
    adminEmailElement.textContent = email;
  }
  if (settingsEmailElement) {
    settingsEmailElement.textContent = email;
  }
}
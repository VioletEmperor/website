// Firebase functions will be loaded from window object
let auth, onAuthStateChanged, signOut;
// Auth state management
let currentUser = null;
let authToken = null;

// DOM elements
let adminEmailElement, settingsEmailElement, logoutButton;
let navTabs, tabContents, fileUpload, fileInput, fileInfo;
let editFileUpload, editFileInput, editFileInfo, editTab;
let currentEditPostId = null;

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
  // Initialize edit tab elements
  editTab = document.querySelector('[data-tab="edit"]');
  editFileUpload = document.getElementById('edit-file-upload');
  editFileInput = document.getElementById('edit-file-input');
  editFileInfo = document.getElementById('edit-file-info');

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

  // Edit file upload functionality
  if (editFileUpload && editFileInput) {
    setupEditFileUpload();
  }

  // Delete confirmation
  const deleteButtons = document.querySelectorAll('.btn-delete');
  deleteButtons.forEach(button => {
    button.addEventListener('click', handleDeletePost);
  });

  // Edit post functionality
  const editButtons = document.querySelectorAll('.btn-edit');
  editButtons.forEach(button => {
    button.addEventListener('click', handleEditPost);
  });

  // Logout functionality
  if (logoutButton) {
    logoutButton.addEventListener('click', handleLogout);
  }

  // Form submissions
  const uploadForm = document.querySelector('.upload-form');
  if (uploadForm) {
    uploadForm.addEventListener('submit', handleUploadFormSubmit);
  }

  const editForm = document.querySelector('#edit .upload-form');
  if (editForm) {
    editForm.addEventListener('submit', handleEditFormSubmit);
  }

  // Cancel edit button
  const cancelEditButton = document.getElementById('cancel-edit');
  if (cancelEditButton) {
    cancelEditButton.addEventListener('click', cancelEdit);
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

function setupEditFileUpload() {
  // Click to upload
  editFileUpload.addEventListener('click', () => {
    editFileInput.click();
  });

  // Drag and drop
  editFileUpload.addEventListener('dragover', (e) => {
    e.preventDefault();
    editFileUpload.classList.add('dragover');
  });

  editFileUpload.addEventListener('dragleave', () => {
    editFileUpload.classList.remove('dragover');
  });

  editFileUpload.addEventListener('drop', (e) => {
    e.preventDefault();
    editFileUpload.classList.remove('dragover');
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      editFileInput.files = files;
      showEditFileInfo(files[0]);
    }
  });

  // File selection
  editFileInput.addEventListener('change', () => {
    if (editFileInput.files.length > 0) {
      showEditFileInfo(editFileInput.files[0]);
    }
  });
}

function showEditFileInfo(file) {
  if (editFileInfo) {
    editFileInfo.textContent = `Selected: ${file.name} (${(file.size / 1024).toFixed(1)} KB)`;
    editFileInfo.style.display = 'block';
  }
}

async function handleDeletePost(e) {
  e.preventDefault();
  const postId = e.target.getAttribute('data-id');
  const postTitle = e.target.getAttribute('data-title') || 'this post';
  
  if (!postId) {
    alert('Post ID not found');
    return;
  }
  
  if (confirm(`Are you sure you want to delete "${postTitle}"? This action cannot be undone.`)) {
    try {
      const response = await fetch(`/admin/posts/${postId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json'
        }
      });
      
      if (response.ok) {
        // Remove the row from the table
        e.target.closest('tr').remove();
        alert('Post deleted successfully');
      } else {
        const error = await response.text();
        alert(`Failed to delete post: ${error}`);
      }
    } catch (error) {
      console.error('Delete error:', error);
      alert('Failed to delete post: Network error');
    }
  }
}

async function handleEditPost(e) {
  e.preventDefault();
  const postId = e.target.getAttribute('data-id');
  
  if (!postId) {
    alert('Post ID not found');
    return;
  }
  
  try {
    // Get post data
    const response = await fetch(`/admin/posts/${postId}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch post data');
    }
    
    const post = await response.json();
    
    // Store current edit post ID
    currentEditPostId = postId;
    
    // Show and switch to edit tab
    if (editTab) {
      editTab.style.display = 'block';
    }
    switchTab('edit');
    
    // Populate the edit form with post data
    const editTitleInput = document.getElementById('edit-title');
    const editExcerptInput = document.getElementById('edit-excerpt');
    const editPostTitle = document.getElementById('edit-post-title');
    
    if (editTitleInput) editTitleInput.value = post.Title || '';
    if (editExcerptInput) editExcerptInput.value = post.Description || '';
    if (editPostTitle) editPostTitle.textContent = post.Title || 'Unknown';
    
    // Clear any previously selected file
    if (editFileInput) {
      editFileInput.value = '';
    }
    if (editFileInfo) {
      editFileInfo.style.display = 'none';
      editFileInfo.textContent = '';
    }
    
  } catch (error) {
    console.error('Edit error:', error);
    alert('Failed to load post for editing');
  }
}

function cancelEdit() {
  // Reset edit form
  const editTitleInput = document.getElementById('edit-title');
  const editExcerptInput = document.getElementById('edit-excerpt');
  const editPostTitle = document.getElementById('edit-post-title');
  
  if (editTitleInput) editTitleInput.value = '';
  if (editExcerptInput) editExcerptInput.value = '';
  if (editPostTitle) editPostTitle.textContent = 'Loading...';
  
  // Clear file input
  if (editFileInput) {
    editFileInput.value = '';
  }
  if (editFileInfo) {
    editFileInfo.style.display = 'none';
    editFileInfo.textContent = '';
  }
  
  // Hide edit tab and switch back to posts
  if (editTab) {
    editTab.style.display = 'none';
  }
  currentEditPostId = null;
  switchTab('posts');
}

async function handleLogout(e) {
  e.preventDefault();
  
  if (confirm('Are you sure you want to logout?')) {
    await signOutUser();
  }
}

async function handleUploadFormSubmit(e) {
  e.preventDefault();
  
  // Create FormData from the form
  const formData = new FormData(e.target);
  
  // Add auth token to form data
  if (authToken) {
    formData.append('authToken', authToken);
  }
  
  try {
    const response = await fetch('/admin/posts/upload', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authToken}`
      },
      body: formData
    });
    
    if (response.ok) {
      alert('Post uploaded successfully');
      // Reset form
      e.target.reset();
      if (fileInfo) {
        fileInfo.style.display = 'none';
        fileInfo.textContent = '';
      }
      // Reload to show new post
      location.reload();
    } else {
      const error = await response.text();
      alert(`Failed to upload post: ${error}`);
    }
  } catch (error) {
    console.error('Upload error:', error);
    alert('Failed to upload post: Network error');
  }
}

async function handleEditFormSubmit(e) {
  e.preventDefault();
  
  if (!currentEditPostId) {
    alert('No post selected for editing');
    return;
  }
  
  const titleInput = document.getElementById('edit-title');
  const excerptInput = document.getElementById('edit-excerpt');
  const fileInput = document.getElementById('edit-file-input');
  
  const hasNewFile = fileInput && fileInput.files && fileInput.files.length > 0;
  
  if (hasNewFile) {
    // Handle file upload - use multipart form data
    const formData = new FormData();
    formData.append('title', titleInput.value);
    formData.append('excerpt', excerptInput.value);
    formData.append('htmlFile', fileInput.files[0]);
    formData.append('editMode', 'true');
    formData.append('postId', currentEditPostId);
    formData.append('authToken', authToken);
    
    try {
      const response = await fetch('/admin/posts/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authToken}`
        },
        body: formData
      });
      
      if (response.ok) {
        alert('Post updated successfully');
        cancelEdit();
        location.reload();
      } else {
        const error = await response.text();
        alert(`Failed to update post: ${error}`);
      }
    } catch (error) {
      console.error('Update error:', error);
      alert('Failed to update post: Network error');
    }
  } else {
    // Handle metadata-only update - use JSON
    const updatedData = {
      title: titleInput.value,
      description: excerptInput.value
    };
    
    try {
      const response = await fetch(`/admin/posts/${currentEditPostId}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updatedData)
      });
      
      if (response.ok) {
        alert('Post updated successfully');
        cancelEdit();
        location.reload();
      } else {
        const error = await response.text();
        alert(`Failed to update post: ${error}`);
      }
    } catch (error) {
      console.error('Update error:', error);
      alert('Failed to update post: Network error');
    }
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
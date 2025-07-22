// Firebase functions will be loaded from window object
let auth, signInWithEmailAndPassword, signInWithPopup, GoogleAuthProvider, onAuthStateChanged;

// Auth state management
let isAuthenticated = false;
let currentUser = null;

// DOM elements
let emailForm, emailInput, passwordInput, gmailButton, errorContainer, tabButtons, tabContents;

document.addEventListener('DOMContentLoaded', function() {
  // Initialize DOM elements
  emailForm = document.querySelector('#email-login .login-form');
  emailInput = document.getElementById('email');
  passwordInput = document.getElementById('password');
  gmailButton = document.querySelector('.btn-gmail');
  errorContainer = document.querySelector('.error-message');
  tabButtons = document.querySelectorAll('.tab-button');
  tabContents = document.querySelectorAll('.tab-content');

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
  signInWithEmailAndPassword = window.signInWithEmailAndPassword;
  signInWithPopup = window.signInWithPopup;
  GoogleAuthProvider = window.GoogleAuthProvider;
  onAuthStateChanged = window.onAuthStateChanged;

  // Initialize authentication
  initializeAuth();
  
  // Set up event listeners
  setupEventListeners();
  
  // Check if user is already authenticated
  checkAuthState();
}

function initializeAuth() {
  // Listen for auth state changes
  onAuthStateChanged(auth, async (user) => {
    if (user) {
      currentUser = user;
      isAuthenticated = true;
      
      // Verify user is admin
      if (await verifyAdminAccess(user.email)) {
        redirectToDashboard();
      } else {
        showError('You are not authorized to access the admin area.');
        await signOut();
      }
    } else {
      currentUser = null;
      isAuthenticated = false;
    }
  });
}

function setupEventListeners() {
  // Tab switching
  tabButtons.forEach(button => {
    button.addEventListener('click', () => {
      const targetTab = button.getAttribute('data-tab');
      switchTab(targetTab);
    });
  });

  // Email/password form submission
  if (emailForm) {
    emailForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      await handleEmailLogin();
    });
  }

  // Gmail sign-in button
  if (gmailButton) {
    gmailButton.addEventListener('click', async (e) => {
      e.preventDefault();
      await handleGmailLogin();
    });
  }
}

function switchTab(targetTab) {
  // Remove active class from all tabs and contents
  tabButtons.forEach(btn => btn.classList.remove('active'));
  tabContents.forEach(content => content.classList.remove('active'));

  // Add active class to clicked tab and corresponding content
  document.querySelector(`[data-tab="${targetTab}"]`).classList.add('active');
  document.getElementById(targetTab + '-login').classList.add('active');
  
  // Clear any errors when switching tabs
  clearError();
}

async function handleEmailLogin() {
  const email = emailInput.value.trim();
  const password = passwordInput.value.trim();

  if (!email || !password) {
    showError('Please enter both email and password.');
    return;
  }

  try {
    showLoading(true);
    clearError();
    
    // Sign in with Firebase
    await signInWithEmailAndPassword(auth, email, password);
    
    // Auth state change will handle the rest
  } catch (error) {
    console.error('Email login error:', error);
    showError(getErrorMessage(error));
  } finally {
    showLoading(false);
  }
}

async function handleGmailLogin() {
  try {
    showLoading(true);
    clearError();
    
    const provider = new GoogleAuthProvider();
    provider.addScope('email');
    provider.addScope('profile');
    
    await signInWithPopup(auth, provider);
    
    // Auth state change will handle the rest
  } catch (error) {
    console.error('Gmail login error:', error);
    showError(getErrorMessage(error));
  } finally {
    showLoading(false);
  }
}

async function verifyAdminAccess(email) {
  try {
    // Get ID token
    const idToken = await currentUser.getIdToken();
    
    // Verify with backend
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

async function signOut() {
  try {
    await window.signOut(auth);
    // Clear any stored tokens
    localStorage.removeItem('adminToken');
    // Clear cookie
    document.cookie = 'adminToken=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
  } catch (error) {
    console.error('Sign out error:', error);
  }
}

function redirectToDashboard() {
  // Store auth token for API calls
  currentUser.getIdToken().then(token => {
    localStorage.setItem('adminToken', token);
    // Set cookie for server-side middleware
    document.cookie = `adminToken=${token}; path=/; secure; samesite=strict`;
    window.location.href = '/admin/dashboard';
  });
}

function checkAuthState() {
  // Check if user is already signed in
  const user = auth.currentUser;
  if (user) {
    // User is already authenticated, verify admin access
    verifyAdminAccess(user.email).then(isAdmin => {
      if (isAdmin) {
        redirectToDashboard();
      }
    });
  }
}

function showError(message) {
  if (errorContainer) {
    errorContainer.textContent = message;
    errorContainer.style.display = 'block';
  }
}

function clearError() {
  if (errorContainer) {
    errorContainer.style.display = 'none';
  }
}

function showLoading(show) {
  const submitButton = document.querySelector('#email-login .btn-primary');
  const gmailBtn = document.querySelector('.btn-gmail');
  
  if (show) {
    if (submitButton) submitButton.textContent = 'Signing in...';
    if (gmailBtn) gmailBtn.textContent = 'Signing in...';
  } else {
    if (submitButton) submitButton.textContent = 'Sign In';
    if (gmailBtn) gmailBtn.innerHTML = `
      <svg class="gmail-icon" viewBox="0 0 24 24" width="20" height="20">
        <path fill="#4285f4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
        <path fill="#34a853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
        <path fill="#fbbc05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
        <path fill="#ea4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
      </svg>
      Sign in with Gmail
    `;
  }
}

function getErrorMessage(error) {
  switch (error.code) {
    case 'auth/invalid-email':
      return 'Invalid email address.';
    case 'auth/user-disabled':
      return 'This account has been disabled.';
    case 'auth/user-not-found':
      return 'No account found with this email.';
    case 'auth/wrong-password':
      return 'Incorrect password.';
    case 'auth/invalid-credential':
      return 'Invalid email or password.';
    case 'auth/too-many-requests':
      return 'Too many failed attempts. Please try again later.';
    case 'auth/popup-closed-by-user':
      return 'Sign-in popup was closed. Please try again.';
    case 'auth/popup-blocked':
      return 'Sign-in popup was blocked. Please allow popups and try again.';
    case 'auth/cancelled-popup-request':
      return 'Sign-in was cancelled. Please try again.';
    default:
      return 'An error occurred during sign-in. Please try again.';
  }
}
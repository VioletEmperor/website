let currentSlideIndex = 0;
const slides = document.querySelectorAll('.carousel-item');
const dots = document.querySelectorAll('.dot');
let isAnimating = false;

function initializeCarousel() {
  // Hide all slides except the first one
  slides.forEach((slide, index) => {
    slide.classList.remove('active', 'slide-out-left', 'slide-in-right');
    if (index === 0) {
      slide.classList.add('active');
    } else {
      slide.style.transform = 'translateX(100%)';
    }
  });
  
  // Set first dot as active
  if (dots.length > 0) {
    dots[0].classList.add('active');
  }
}

function showSlide(index, direction = 1) {
  if (isAnimating || index === currentSlideIndex) return;
  
  isAnimating = true;
  const currentSlide = slides[currentSlideIndex];
  const nextSlide = slides[index];
  
  // Clear all animation classes
  slides.forEach(slide => {
    slide.classList.remove('slide-out-left', 'slide-in-right');
  });
  
  if (direction > 0) {
    // Moving forward: current slide goes left, next slide comes from right
    currentSlide.classList.remove('active');
    currentSlide.classList.add('slide-out-left');
    
    nextSlide.classList.add('slide-in-right');
    setTimeout(() => {
      nextSlide.classList.remove('slide-in-right');
      nextSlide.classList.add('active');
    }, 50);
  } else {
    // Moving backward: current slide goes right, next slide comes from left
    currentSlide.classList.remove('active');
    currentSlide.style.transform = 'translateX(100%)';
    
    nextSlide.style.transform = 'translateX(-100%)';
    setTimeout(() => {
      nextSlide.style.transform = 'translateX(0)';
      nextSlide.classList.add('active');
    }, 50);
  }
  
  // Update dots
  dots.forEach(dot => dot.classList.remove('active'));
  if (dots[index]) {
    dots[index].classList.add('active');
  }
  
  // Clean up after animation
  setTimeout(() => {
    slides.forEach(slide => {
      if (!slide.classList.contains('active')) {
        slide.classList.remove('slide-out-left', 'slide-in-right');
        slide.style.transform = 'translateX(100%)';
      }
    });
    isAnimating = false;
  }, 650);
}

function changeSlide(direction) {
  const newIndex = (currentSlideIndex + direction + slides.length) % slides.length;
  currentSlideIndex = newIndex;
  showSlide(currentSlideIndex, direction);
}

function currentSlide(index) {
  const direction = index - 1 > currentSlideIndex ? 1 : -1;
  currentSlideIndex = index - 1;
  showSlide(currentSlideIndex, direction);
}

// Initialize carousel when DOM is loaded
document.addEventListener('DOMContentLoaded', initializeCarousel);

// Auto-advance carousel every 5 seconds
setInterval(() => {
  changeSlide(1);
}, 5000);
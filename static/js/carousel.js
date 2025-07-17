let currentSlideIndex = 0;
let isTransitioning = false;
const slides = document.querySelectorAll('.carousel-item');
const dots = document.querySelectorAll('.dot');

function initializeCarousel() {
  if (slides.length > 0) {
    slides[0].classList.add('active');
    dots[0].classList.add('active');
  }
  setInterval(nextSlide, 10000);
}

function nextSlide() {
  showSlide(Math.abs((currentSlideIndex + 1) % slides.length));
}

function prevSlide() {
  showSlide(Math.abs((currentSlideIndex - 1) % slides.length));
}

function showSlide(index) {
  if (isTransitioning) return;
  
  isTransitioning = true;
  currentSlideIndex = index;
  
  slides.forEach(slide => {
    slide.classList.remove("active");
  });
  
  dots.forEach(dot => {
    dot.classList.remove("active");
  });

  slides[currentSlideIndex].classList.add("active");
  dots[currentSlideIndex].classList.add("active");
  
  setTimeout(() => {
    isTransitioning = false;
  }, 600);
}

document.addEventListener('DOMContentLoaded', initializeCarousel);
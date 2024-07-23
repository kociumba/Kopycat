let isGayMode = false;
let animationId;

function rainbowAnimation() {
  let hue = Math.floor(Math.random() * 360);

  function updateColor() {
    hue = (hue + 1) % 360;
    document.documentElement.style.setProperty('--primary-color', `hsl(${hue}, 100%, 50%)`);
    animationId = requestAnimationFrame(updateColor);
  }

  updateColor();
}

function stopAnimation() {
  if (animationId) {
    cancelAnimationFrame(animationId);
    animationId = null;
  }
}

function toggleGayMode() {
  isGayMode = !isGayMode;
  if (isGayMode) {
    rainbowAnimation();
  } else {
    stopAnimation();
    document.documentElement.style.setProperty('--primary-color', '#13995a');
  }
}

// Add click event listener to the entire document
$(document).on('click', function(e) {
  // Check if the clicked element or any of its parents have the id 'gay-mode-trigger'
  if ($(e.target).closest('#gay-toggle').length) {
    toggleGayMode();
  }
});

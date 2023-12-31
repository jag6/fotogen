//MOBILE NAV
const mobile_nav = document.getElementById('mobile-nav');
const nav_overlay = document.getElementById('nav-overlay');
const body = document.querySelector('body');

const clickOff = () => {
    nav_overlay.style.display = 'none';
    mobile_nav.style.width = '0%';
    body.style.overflow = 'auto';
};
document.getElementById('hamburger-icon').addEventListener('click', () => {
    if(mobile_nav.style.width === '300px') {
        clickOff();
    }else {
        mobile_nav.style.width = '300px';
        nav_overlay.style.display = 'flex';
        body.style.overflow = 'hidden';
    }
});
nav_overlay.addEventListener('click', (e) => {
    if(e.target == nav_overlay) {
        clickOff();
    }
});


//ANIMATIONS
const scrollElements = document.querySelectorAll('.js-scroll');
var throttleTimer;
const throttle = (callback, time) => {
    if (throttleTimer) return;
    throttleTimer = true;
    setTimeout(() => {
        callback();
        throttleTimer = false;
        }, time
    );
}
const elementInView = (el, dividend = 1) => {
const elementTop = el.getBoundingClientRect().top;
	return (
    	elementTop <= (window.innerHeight || document.documentElement.clientHeight) / dividend
	);
};
const elementOutofView = (el) => {
	const elementTop = el.getBoundingClientRect().top;
	return (
    	elementTop > (window.innerHeight || document.documentElement.clientHeight)
  );
};
const displayScrollElement = (element) => {
  	element.classList.add('scrolled');
};
const hideScrollElement = (element) => {
  	element.classList.remove('scrolled');
};
const handleScrollAnimation = () => {
  	scrollElements.forEach((el) => {
    	if (elementInView(el, 1.25)) {
      	displayScrollElement(el);
    	}else if (elementOutofView(el)) {
      	hideScrollElement(el)
    }
  })
};
window.addEventListener("scroll", () => { 
    throttle(() => {
        handleScrollAnimation();
    }, 250);
});


//ALERTS
document.querySelectorAll('.alert-message').forEach((alert) => {
	//auto delete
	setTimeout(() => {
		alert.remove();
	}, 5000);
	//click to delete
	alert.addEventListener('click', () => {
		alert.classList.add('fade-out');
		setTimeout(() => {
			alert.remove();
		}, 300);
	});
});


//RENDER
// document.querySelector('.render').innerHTML = `
   
// `;
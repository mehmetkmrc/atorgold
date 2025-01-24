// JavaScript: Side Panel Toggle

let jsAccessEnabled = false;

document.addEventListener('DOMContentLoaded', () => {
    // Function to toggle side panel configuration
    function toggleSidePanel(position) {
        if (!jsAccessEnabled) {
            console.warn('Access to this functionality is restricted.');
            return;
        }
        

        const body = document.body;
        

        // Reset classes to avoid conflicts
        body.className = 'stretched';
        

        if (position === 'left') {
            body.classList.add(
                'side-panel-left',
                'device-down-xxl',
                'device-lg',
                'device-down-xl',
                'device-up-lg',
                'device-up-md',
                'device-up-sm',
                'device-up-xs',
                'is-expanded-menu',
                'has-plugin-onepage',
                'has-plugin-bootstrap',
                'has-plugin-counter',
                'has-plugin-navtree',
                'has-plugin-form',
                'quick-contact-form-ready',
                'vsc-initialized',
                'has-plugin-tips',
                'init-plugin-tips',
                'has-plugin-notify',
                'side-panel-open',
                'has-plugin-html5video',
            );
            
        } else if (position === 'right') {
            body.classList.add(
                'has-plugin-html5video',
                'device-down-xxl',
                'device-lg',
                'device-down-xl',
                'device-up-lg',
                'device-up-md',
                'device-up-sm',
                'device-up-xs',
                'is-expanded-menu',
                'has-plugin-onepage',
                'has-plugin-bootstrap',
                'has-plugin-counter',
                'has-plugin-navtree',
                'has-plugin-form',
                'quick-contact-form-ready',
                'vsc-initialized',
                'has-plugin-tips',
                'init-plugin-tips',
                'side-panel-open',
                'has-plugin-notify'
            );
           
        }
    }

    // Add click event listeners for toggle buttons
    const leftPanelTrigger = document.querySelector('#left-panel-trigger');
    const rightPanelTrigger = document.querySelector('#right-panel-trigger');

    if (leftPanelTrigger) {
        leftPanelTrigger.addEventListener('click', (e) => {
            e.preventDefault();
            toggleSidePanel('left');
        });
    }

    if (rightPanelTrigger) {
        rightPanelTrigger.addEventListener('click', (e) => {
            e.preventDefault();
            toggleSidePanel('right');
        });
    }

    // Enable access when specific button is clicked
    const enableAccessButton = document.querySelector('#enable-js-access');
    if (enableAccessButton) {
        enableAccessButton.addEventListener('click', (e) => {
            e.preventDefault();
            jsAccessEnabled = true;
            console.log('JavaScript access has been enabled.');
        });
    }
});

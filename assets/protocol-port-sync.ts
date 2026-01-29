export function initProtocolPortSync() {
    const protocolSelect = document.querySelector('select[name="protocol"]') as HTMLSelectElement | null;
    const portInput = document.querySelector('input[name="expose_port"]') as HTMLInputElement | null;

    if (!protocolSelect || !portInput) return;

    protocolSelect.addEventListener('change', function () {
        if (this.value === 'https') {
            portInput.value = '443';
        } else if (this.value === 'http') {
            portInput.value = '80';
        }
    });
}

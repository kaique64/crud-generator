/**
 * main.js
 * * Lógica de frontend aprimorada com IMask.js e validação real-time.
 */

// Espera o DOM estar pronto
document.addEventListener('DOMContentLoaded', () => {

    // --- Seletores de Elementos ---
    const form = document.getElementById('crud-form');
    if (!form) return; // Sai se o formulário não existir

    const formTitle = document.getElementById('form-title');
    const formSubmitBtn = document.getElementById('form-submit-btn');
    const formCancelBtn = document.getElementById('form-cancel-btn');
    const formIdField = document.getElementById('form-id-field');
    const formCard = document.getElementById('form-card');
    const formInputs = form.querySelectorAll('input[name]');

    // --- Estado do Formulário ---
    const originalFormAction = form.action;
    const originalFormTitle = formTitle.innerText;
    const originalSubmitText = formSubmitBtn.innerText;
    let inputMasks = []; // Armazena instâncias do IMask

    // --- 1. INICIALIZAÇÃO ---

    /**
     * Inicializa as máscaras em todos os inputs com 'data-mask'
     */
     const initMasks = () => {
       const customDefinitions = {
            // '9': Dígito (já definido implicitamente na maioria das vezes,
            // mas vamos garantir que ele seja mapeado para 0-9)
            '9': {
                mask: /\d/,
                lazy: false
            },
            // '#': Caractere (letra A-Z, a-z)
            '#': {
                mask: /[a-zA-Z]/,
                lazy: false
            },
            // '*': Qualquer Caractere/Tipo (inclui dígitos, letras e outros)
            '*': {
                mask: /./, // Aceita qualquer caractere
                lazy: false
            }
        };
        inputMasks = []; // Limpa máscaras antigas
        formInputs.forEach(input => {
            const maskPattern = input.dataset.mask; // ex: "999.999.999-99"
            if (!maskPattern) return;

            // --- INÍCIO DA CORREÇÃO ---
            // IMask.js usa '0' como placeholder de dígito, não '9'.
            // Vamos converter nosso padrão "9" para o padrão "0" do IMask.
            const imaskPattern = maskPattern.replace(/9/g, '0');
            // --- FIM DA CORREÇÃO ---

            let maskOptions = {
                mask: imaskPattern, // <-- USA A VARIÁVEL CORRIGIDA
                lazy: false, // Preenche a máscara à medida que o usuário digita
                definitions: customDefinitions
            };

            // Lógica especial para máscaras dinâmicas (ex: Telefone)
            // O schema original é "(99) 99999-9999" (com 9)
            // O IMask precisa de "(00) 00000-0000" (com 0)
            if (maskPattern.includes('(99) 9')) {
                maskOptions.mask = [
                    {
                        mask: '(00) 0000-0000', // Convertido para '0'
                        lazy: false
                    },
                    {
                        mask: '(00) 00000-0000', // Convertido para '0'
                        lazy: false
                    }
                ];
            }

            // Cria a instância do IMask e armazena
            const maskInstance = IMask(input, maskOptions);
            inputMasks.push(maskInstance);
        });
    };

    /**
     * Inicializa os listeners de validação
     */
    const initValidation = () => {
        formInputs.forEach(input => {
            // Valida quando o usuário *sai* do campo
            input.addEventListener('blur', (e) => {
                validateField(e.target);
            });

            // Limpa o erro assim que o usuário começa a corrigir
            input.addEventListener('input', (e) => {
                clearError(e.target);
            });
        });
    };

    // --- 2. LÓGICA DE VALIDAÇÃO ---

    /**
     * Valida um único campo e exibe/oculta a mensagem de erro
     * @param {HTMLInputElement} input
     * @returns {boolean} - True se for válido, False se for inválido
     */
    const validateField = (input) => {
        const value = input.value;
        const type = input.dataset.validateType;
        const isRequired = input.hasAttribute('required');

        // 1. Validação de Obrigatório
        if (isRequired && value.trim() === '') {
            showError(input, 'Campo obrigatório');
            return false;
        }

        // 2. Validações de Tipo (CPF, Email, etc.)
        if (value.trim() !== '') { // Só valida o formato se não estiver vazio
            switch (type) {
                case 'cpf':
                    if (!isValidCPF(value, isRequired)) {
                        showError(input, 'CPF inválido');
                        return false;
                    }
                    break;
                case 'email':
                    if (!isValidEmail(value, isRequired)) {
                        showError(input, 'Email inválido');
                        return false;
                    }
                    break;
                case 'cnpj':
                    if (!isValidCNPJ(value, isRequired)) {
                        showError(input, 'CNPJ inválido');
                        return false;
                    }
                    break;
                case 'telefone':
                    if (!isValidTelefone(value, isRequired)) {
                        showError(input, 'Telefone inválido');
                        return false;
                    }
                    break;
                case 'cep':
                    if (!isValidCEP(value, isRequired)) {
                        showError(input, 'CEP inválido');
                        return false;
                    }
                    break;
                // Adicione casos para 'cnpj', 'telefone', 'cep' aqui
            }
        }

        // 3. Validações de Regex (se passadas do backend)
        // (Esta parte pode ser adicionada depois, lendo data-regex-pattern)

        // Se passou por tudo, está válido
        clearError(input);
        return true;
    };

    /**
     * Exibe a mensagem de erro para um campo
     * @param {HTMLInputElement} input
     * @param {string} message
     */
    const showError = (input, message) => {
        const errorJS = document.getElementById(`error-js-${input.name}`);
        const errorBackend = document.getElementById(`error-backend-${input.name}`);

        // Oculta erro do backend (se houver) para dar lugar ao erro de JS
        if (errorBackend) errorBackend.classList.add('hidden');

        if (errorJS) {
            errorJS.innerText = message;
            errorJS.classList.remove('hidden');
        }
        input.classList.add('border-red-500', 'ring-1', 'ring-red-500');
    };

    /**
     * Limpa a mensagem de erro de um campo
     * @param {HTMLInputElement} input
     */
    const clearError = (input) => {
        const errorJS = document.getElementById(`error-js-${input.name}`);
        const errorBackend = document.getElementById(`error-backend-${input.name}`);

        if (errorBackend) errorBackend.classList.add('hidden'); // Esconde permanente
        if (errorJS) errorJS.classList.add('hidden');

        input.classList.remove('border-red-500', 'ring-1', 'ring-red-500');
    };

    /**
     * Limpa todos os erros de validação do formulário
     */
    const clearAllValidation = () => {
        formInputs.forEach(input => clearError(input));
    };

    /**
     * Valida o formulário inteiro antes do envio
     * @returns {boolean} - True se o formulário for válido
     */
    const validateForm = () => {
        let isFormValid = true;
        formInputs.forEach(input => {
            // Se 'validateField' retornar false, o formulário é inválido
            if (!validateField(input)) {
                isFormValid = false;
            }
        });
        return isFormValid;
    };

    // --- 3. LÓGICA DE EDIÇÃO / CRIAÇÃO ---

    /**
     * Prepara o formulário para edição (chamado pelo HTML)
     * @param {string} id
     */
    window.startEdit = async (id) => {
        try {
            const response = await fetch(`/get?id=${id}`);
            if (!response.ok) throw new Error('Falha ao carregar dados');

            const data = await response.json();

            // Reseta o formulário para modo "limpo"
            cancelEdit();

            // Popula os campos
            formInputs.forEach(input => {
                if (data[input.name]) {
                    let value = data[input.name];

                    // Trata datas
                    if (input.type === 'date' && value) {
                        value = value.split('T')[0]; // Formato AAAA-MM-DD
                    }

                    input.value = value;

                    // IMPORTANTE: Atualiza o valor da máscara
                    const mask = inputMasks.find(m => m.el === input);
                    if (mask) {
                        mask.updateValue();
                    }
                }
            });

            // Atualiza UI do formulário para modo "Edição"
            formIdField.value = id;
            form.action = `/update?id=${id}`;
            formTitle.innerText = `Editando Registro #${id}`;
            formSubmitBtn.innerText = 'Atualizar';
            formCancelBtn.style.display = 'inline-block';

            formCard.scrollIntoView({ behavior: 'smooth' });

        } catch (error) {
            console.error('Falha ao buscar dados para edição:', error);
            alert('Não foi possível carregar os dados para edição.');
        }
    };

    /**
     * Reseta o formulário para o modo "Criação"
     */
    const cancelEdit = () => {
        form.reset();
        clearAllValidation();

        // Reseta os valores das máscaras
        inputMasks.forEach(mask => mask.updateValue());

        formIdField.value = '';
        form.action = originalFormAction;
        formTitle.innerText = originalFormTitle;
        formSubmitBtn.innerText = originalSubmitText;
        formCancelBtn.style.display = 'none';
    };

    // --- 4. EVENT LISTENERS ---

    // Listener do botão Cancelar
    formCancelBtn.addEventListener('click', cancelEdit);

    // Listener de Submissão do Formulário
    form.addEventListener('submit', (e) => {
        // Previne o envio se a validação de frontend falhar
        if (!validateForm()) {
            e.preventDefault();
            console.warn('Formulário inválido. Verifique os campos.');
            // Foca no primeiro campo inválido
            form.querySelector('.is-invalid')?.focus();
        }
        // Se validateForm() for true, o formulário é enviado normalmente
    });

    // --- 5. HABILIDADES (Funções de Validação) ---
    // (Estas são duplicatas da lógica do backend para UX)

    const isValidEmail = (email, isRequired) => {
        if (email === '' && !isRequired) return true;
      
        return /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(email);
    };

    const isValidCPF = (cpf, isRequired) => {      
        cpf = cpf.replace(/\D/g, ''); // Remove não-dígitos

        if (cpf === '' && !isRequired) return true;
        
        if (cpf.length !== 11 || /^(\d)\1{10}$/.test(cpf)) return false;

        let sum = 0, rest;

        for (let i = 1; i <= 9; i++) sum += parseInt(cpf.substring(i-1, i)) * (11 - i);
        rest = (sum * 10) % 11;
        if ((rest === 10) || (rest === 11)) rest = 0;
        if (rest !== parseInt(cpf.substring(9, 10))) return false;

        sum = 0;
        for (let i = 1; i <= 10; i++) sum += parseInt(cpf.substring(i-1, i)) * (12 - i);
        rest = (sum * 10) % 11;
        if ((rest === 10) || (rest === 11)) rest = 0;
        if (rest !== parseInt(cpf.substring(10, 11))) return false;

        return true;
    };

    const isValidCNPJ = (cnpj, isRequired) => {
      cnpj = cnpj.replace(/[^\d]+/g, ''); // Remove caracteres não numéricos
      
      if (cnpj === '' && !isRequired) return true;
      
      if (cnpj.length !== 14) return false; // Verifica se o CNPJ tem 14 dígitos
  
      // Validação de CNPJ com números repetidos (ex: 11111111111111)
      if (/^(.)\1{13}$/.test(cnpj)) return false;
  
      // Validação do primeiro dígito verificador
      let soma = 0;
      let peso = [6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2];
      for (let i = 0; i < 12; i++) {
          soma += parseInt(cnpj[i]) * peso[i + 1];
      }
      let digito1 = 11 - (soma % 11);
      digito1 = digito1 === 10 || digito1 === 11 ? 0 : digito1;
  
      // Validação do segundo dígito verificador
      soma = 0;
      peso = [5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2];
      for (let i = 0; i < 13; i++) {
          soma += parseInt(cnpj[i]) * peso[i];
      }
      let digito2 = 11 - (soma % 11);
      digito2 = digito2 === 10 || digito2 === 11 ? 0 : digito2;
  
      // Verifica se os dígitos calculados são iguais aos do CNPJ
      return cnpj[12] == digito1 && cnpj[13] == digito2;
    };
    
    // NOVO: Validação de Telefone
    const isValidTelefone = (telefone, isRequired) => {        
        telefone = telefone.replace(/\D/g, ''); // Remove não-dígitos
        
        if (telefone === '' && !isRequired) return true;
        
        // Aceita 10 (DDD + 8 dígitos) ou 11 (DDD + 9 + 8 dígitos)
        return telefone.length >= 10 && telefone.length <= 11;
    };
    
    // NOVO: Validação de CEP
    const isValidCEP = (cep, isRequired) => {
        cep = cep.replace(/\D/g, ''); // Remove não-dígitos  

        if (cep === '' && !isRequired) return true;

        // O CEP sempre deve ter 8 dígitos
        return cep.length === 8;
    };
    

    // --- INICIA A MÁGICA ---
    initMasks();
    initValidation();
});

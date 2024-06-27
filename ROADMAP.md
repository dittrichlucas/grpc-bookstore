## ROADMAP

- Adicionar novas rotas
- Refatorar o código para tornar o comando de login utilizável
- Salvar o token gerado pelo login em um arquivo que deverá ser consumido pelo front (?)
- O token deve ser gerado dentro de uma pasta chamada cache.
    - Cada sessão deve ser salva em um arquivo json o nome sendo um uuid
    - Dentro desse arquivo vai ter o token e a sua expiração

- REFATORAR TODO O CÓDIGO
    - Arrumar a estrutura de pastas
    - Renomear variáveis e funções
- Estruturar a documentação (README e ROADMAP)
- Subir o código para o github

### Estrutura da aplicação

De modo geral, o fluxo deve ser o seguinte (contém passos futuros - session id):
- O usuário loga na aplicação com seu user + password
- O servidor verifica as informações e retorna um session id
- Com esse session id, o servidor gera um arquivo cujo nome será essa informação e conterá um token e sua expiração
- De posse desse session id, o client busca o respectivo token para autenticar a requisição

```sh
// server + client

        .-------------.          | .------------. |         .--------------------.
        |    LOGIN    |--------> | | SESSION ID | | ------->| TOKEN + EXPIRATION | 
        '-------------'          | '------------' |         '--------------------'
                                                                      ^
                                                                      |
        .--------------.                                              |
        |    CLIENT    |----------------------------------------------'
        '--------------'         

```

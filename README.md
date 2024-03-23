# HACKATON - Software Architecture

# Description

# Features

## How To Run Locally



### VSCode - Debug
The launch.json file is already configured for debuging. Just hit F5 and be happy. (Local Infra needed to complete experience)

### Local Infra

You can use a local infra to run your application, just following the steps:

## Architecture

Nós escolhemos a arquitetura baseada em microserviços utilizando eventos para os dois designs (MVP e Fase 2) como os diagramas abaixo demonstram. A escolha foi feita por diversos motivos, dentre eles:

- Escalabilidade, a maioria dos componentes da arquitetura pode escalar acima ou abaixo tanto vertical quanto horizontalmente, automaticamente sob demanda sem gerar impactos no atendimendo da regra de negócio.
- Resiliência, se um microserviço falhar o impacto não será generalizado pois os demais serviços que nao dependem dele, continuarão funcionando.
- Desacoplamento, cada microserviço é independente e não se comunica diretamente com os demais, facilitando a substituição e atualizações individuais sem afetar o sistema globalmente.
- Disponibilidade, utilizamos um cluster de kubernetes (EKS) o que garante a capacidade do sistema estar operacional e acessível sem interrupções por conta da replicação e escalabilidade automática, probes de vida e prontidão (healthchecks). Além disso, utilizamos balanceamento de carga entre as aplicações, evitando sobrecargas em instancias individuais.
- Segurança, escolhemos na arquitetura isolar ao máximo as aplicações do acesso externo, utilizando VPCs(rede privada), Subnets privadas, e ferramentas auto-gerenciadas de segurança para acesso externo. Além disso, em um nivel inferior, temos análise estática de codigo e relatorios de segurança (OWASP ZAP) para garantir mais um nível de segurança.
- Manutenibilidade, com domínios subdivididos temos a facilidade de alterar partes específicas das aplicações sem afetar as demais partes. Isto aliado a uma boa documentação feita, facilita ainda mais a evolução do software.
#define sim 2
#define nao 3
#define nao_sei 4

void setup() {
  // put your setup code here, to run once:
  pinMode(sim, INPUT);
  pinMode(nao, INPUT);
  pinMode(nao_sei, INPUT);
  Serial.begin(9600);
}

void loop() {
  // Defino que o arduino deve ler estas conexões com os valores digitais, ou seja, 1 (ligado), 0 (desligo)
  int botao_sim = digitalRead(sim);
  int botao_nao = digitalRead(nao);
  int botao_nao_sei = digitalRead(nao_sei);
  if (botao_sim == HIGH){ // Se o botão for clicado, envia para a porta serial o seguinte valor
    Serial.println("sim"); // Valor a ser enviado
    delay(1000);
  }

  if (botao_nao == HIGH){ // Se o botão for clicado, envia para a porta serial o seguinte valor
    Serial.println("nao"); // Valor a ser enviado
    delay(1000);
  }
  
  if (botao_nao_sei == HIGH){ // Se o botão for clicado, envia para a porta serial o seguinte valor
    Serial.println("nao_sei"); // Valor a ser enviado
    delay(1000);
  }
}

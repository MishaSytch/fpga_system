module fpga (
    input wire        clk,
    input wire        reset_n,
    
    output wire ultrasonic_pulse, // выход генератора УЗ
    input wire ultrasonic_echo,   // вход с приемника УЗ

    // Avalon-MM Slave интерфейс
    input wire [1:0] avl_address,
    input wire avl_read,
    input wire avl_write,
    input wire [31:0] avl_writedata,
    output wire [31:0] avl_readdata
);

// Объявляем сигналы для FIFO
wire write_enable;
wire read_enable;
wire [15:0] echo_data;
wire [15:0] fifo_data_out;

// Генератор ультразвука
ultrasonic_generator u_ultrasonic_gen (
    .clk(clk),
    .rst(~reset_n),          // инвертируем активный низкий сброс
    .pulse_out(ultrasonic_pulse)
);

// Память (FIFO или SRAM)
fifo_memory u_fifo_mem (
    .clk(clk),
    .wr_en(write_enable),
    .data_in(echo_data),
    .rd_en(read_enable),
    .data_out(fifo_data_out)
);

// Детектор фронта сигнала echo
reg echo_d1, echo_d2;
wire echo_rising_edge;

always @(posedge clk or negedge reset_n) begin
    if (!reset_n) begin
        echo_d1 <= 0;
        echo_d2 <= 0;
    end else begin
        echo_d1 <= ultrasonic_echo;
        echo_d2 <= echo_d1;
    end
end

assign echo_rising_edge = (echo_d1 & ~echo_d2); // Поймали фронт

// Управление памятью
assign write_enable = echo_rising_edge;

// Амплитуда сигнала будет равна (например) времени, прошедшему от начала импульса до прихода эха
reg [15:0] echo_amplitude;
always @(posedge clk or negedge reset_n) begin
  if (!reset_n) begin
		echo_amplitude <= 16'd0;
  end else if (echo_rising_edge) begin
		// Сохраняем амплитуду сигнала (например, время, прошедшее с момента импульса)
		echo_amplitude <= 16'hFF; // Здесь 0xFF это просто пример, замените на реальную логику
  end
end

assign echo_data = echo_amplitude;  // Пишем в память амплитуду

// Чтение данных по Avalon-MM
assign read_enable = avl_read & (avl_address == 2'b00);
assign avl_readdata = {16'd0, fifo_data_out};  // Выдаём данные на чтение

endmodule

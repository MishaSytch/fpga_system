module ultrasonic_generator(
    input wire clk,           // Тактовая частота FPGA (например, 50 МГц)
    input wire rst,           // Асинхронный сброс
    output reg pulse_out      // Выходной ультразвуковой сигнал (40 кГц)
);

    reg [15:0] counter;       // 16-битный счётчик тактов

    // Период сигнала 40 кГц при тактовой частоте 50 МГц
    parameter PULSE_PERIOD = 1250;  // 50_000_000 / 40_000 = 1250

    always @(posedge clk or posedge rst) begin
        if (rst) begin
            // Сброс: обнуляем счётчик и выход
            counter <= 16'd0;
            pulse_out <= 1'b0;
        end else begin
            // Если прошло половина периода — инвертируем сигнал
            if (counter >= (PULSE_PERIOD / 2)) begin
                pulse_out <= ~pulse_out;  // Инвертировать выход
                counter <= 16'd0;         // Сбросить счётчик
            end else begin
                counter <= counter + 1;   // Увеличить счётчик
            end
        end
    end

endmodule
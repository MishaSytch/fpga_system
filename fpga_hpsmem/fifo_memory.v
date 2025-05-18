module fifo_memory(
    input wire clk,
    input wire wr_en,
    input wire [15:0] data_in,
    input wire rd_en,
    output reg [15:0] data_out
);
    reg [15:0] fifo [0:15];
    reg [3:0] write_ptr;
    reg [3:0] read_ptr;
	 
	 reg [15:0] count;

    always @(posedge clk) begin
        if (wr_en) begin
				if (count < 1000) begin
					fifo[write_ptr] <= data_in;
					write_ptr <= write_ptr + 1;
					count <= count+1;
				end;
        end
		  else
			count = 0;
    end

    always @(posedge clk) begin
        if (rd_en) begin
            data_out <= fifo[read_ptr];
            read_ptr <= read_ptr + 1;
        end
    end
endmodule

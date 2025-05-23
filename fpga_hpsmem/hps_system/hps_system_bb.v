
module hps_system (
	clk_clk,
	onchip_memory2_0_s1_address,
	onchip_memory2_0_s1_clken,
	onchip_memory2_0_s1_chipselect,
	onchip_memory2_0_s1_write,
	onchip_memory2_0_s1_readdata,
	onchip_memory2_0_s1_writedata,
	onchip_memory2_0_s1_byteenable,
	onchip_memory2_0_reset1_reset,
	onchip_memory2_0_reset1_reset_req);	

	input		clk_clk;
	input	[10:0]	onchip_memory2_0_s1_address;
	input		onchip_memory2_0_s1_clken;
	input		onchip_memory2_0_s1_chipselect;
	input		onchip_memory2_0_s1_write;
	output	[15:0]	onchip_memory2_0_s1_readdata;
	input	[15:0]	onchip_memory2_0_s1_writedata;
	input	[1:0]	onchip_memory2_0_s1_byteenable;
	input		onchip_memory2_0_reset1_reset;
	input		onchip_memory2_0_reset1_reset_req;
endmodule

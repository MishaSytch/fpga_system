	component hps_system is
		port (
			clk_clk                           : in  std_logic                     := 'X';             -- clk
			onchip_memory2_0_s1_address       : in  std_logic_vector(10 downto 0) := (others => 'X'); -- address
			onchip_memory2_0_s1_clken         : in  std_logic                     := 'X';             -- clken
			onchip_memory2_0_s1_chipselect    : in  std_logic                     := 'X';             -- chipselect
			onchip_memory2_0_s1_write         : in  std_logic                     := 'X';             -- write
			onchip_memory2_0_s1_readdata      : out std_logic_vector(15 downto 0);                    -- readdata
			onchip_memory2_0_s1_writedata     : in  std_logic_vector(15 downto 0) := (others => 'X'); -- writedata
			onchip_memory2_0_s1_byteenable    : in  std_logic_vector(1 downto 0)  := (others => 'X'); -- byteenable
			onchip_memory2_0_reset1_reset     : in  std_logic                     := 'X';             -- reset
			onchip_memory2_0_reset1_reset_req : in  std_logic                     := 'X'              -- reset_req
		);
	end component hps_system;

	u0 : component hps_system
		port map (
			clk_clk                           => CONNECTED_TO_clk_clk,                           --                     clk.clk
			onchip_memory2_0_s1_address       => CONNECTED_TO_onchip_memory2_0_s1_address,       --     onchip_memory2_0_s1.address
			onchip_memory2_0_s1_clken         => CONNECTED_TO_onchip_memory2_0_s1_clken,         --                        .clken
			onchip_memory2_0_s1_chipselect    => CONNECTED_TO_onchip_memory2_0_s1_chipselect,    --                        .chipselect
			onchip_memory2_0_s1_write         => CONNECTED_TO_onchip_memory2_0_s1_write,         --                        .write
			onchip_memory2_0_s1_readdata      => CONNECTED_TO_onchip_memory2_0_s1_readdata,      --                        .readdata
			onchip_memory2_0_s1_writedata     => CONNECTED_TO_onchip_memory2_0_s1_writedata,     --                        .writedata
			onchip_memory2_0_s1_byteenable    => CONNECTED_TO_onchip_memory2_0_s1_byteenable,    --                        .byteenable
			onchip_memory2_0_reset1_reset     => CONNECTED_TO_onchip_memory2_0_reset1_reset,     -- onchip_memory2_0_reset1.reset
			onchip_memory2_0_reset1_reset_req => CONNECTED_TO_onchip_memory2_0_reset1_reset_req  --                        .reset_req
		);


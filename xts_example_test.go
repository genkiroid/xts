package xts

var sampleXML = `<?xml version="1.0"?>
<mysqldump xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<database name="test">
	<table_structure name="users">
		<field Field="id" Type="int(11)" Null="NO" Key="PRI" Default="" Extra="auto_increment" Comment="" />
		<field Field="name" Type="varchar(32)" Null="YES" Key="" Extra="" Comment="" />
		<field Field="nick_name" Type="varchar(32)" Null="YES" Key="" Extra="" Comment="" />
		<field Field="state" Type="int(11)" Null="YES" Key="MUL" Extra="" Comment="" />
		<field Field="memo" Type="text" Null="YES" Key="" Extra="" Comment="" />
		<key Table="users" Non_unique="0" Key_name="PRIMARY" Seq_in_index="1" Column_name="id" Collation="A" Cardinality="1" Null="" Index_type="BTREE" Comment="" Index_comment="" />
		<options Name="users" Engine="InnoDB" Version="10" Row_format="Compact" Rows="1" Avg_row_length="99999" Data_length="99999" Max_data_length="0" Index_length="999999" Data_free="9999999999" Create_time="2021-08-09 00:00:00" Collation="ujis_japanese_ci" Create_options="" Comment="" />
	</table_structure>
	<table_data name="users">
	<row>
		<field name="id">1</field>
		<field name="name">Alice</field>
		<field name="nick_name" xsi:nil="true" />
		<field name="state">0</field>
		<field name="memo"></field>
	</row>
	</table_data>
	<table_structure name="orders">
		<field Field="id" Type="int(11)" Null="NO" Key="PRI" Default="" Extra="" Comment="" />
		<field Field="user_id" Type="int(11)" Null="NO" Key="" Extra="" Comment="" />
		<field Field="item_id" Type="int(11)" Null="NO" Key="" Extra="" Comment="" />
		<field Field="order_date" Type="timestamp" Null="NO" Key="" Extra="" Comment="" />
		<key Table="orders" Non_unique="0" Key_name="PRIMARY" Seq_in_index="1" Column_name="id" Collation="A" Cardinality="1" Null="" Index_type="BTREE" Comment="" Index_comment="" />
		<options Name="orders" Engine="InnoDB" Version="10" Row_format="Compact" Rows="1" Avg_row_length="99999" Data_length="99999" Max_data_length="0" Index_length="999999" Data_free="9999999999" Create_time="2021-08-09 00:00:00" Collation="ujis_japanese_ci" Create_options="" Comment="" />
	</table_structure>
	<table_data name="orders">
	<row>
		<field name="id">1</field>
		<field name="user_id">1</field>
		<field name="item_id">1</field>
		<field name="order_date">2021-08-09 12:34:45</field>
	</row>
	<row>
		<field name="id">2</field>
		<field name="user_id">1</field>
		<field name="item_id">2</field>
		<field name="order_date">2021-08-09 13:34:45</field>
	</row>
	</table_data>
</database>
</mysqldump>`

func ExampleSql_InsertStmt() {
	xml, _ := NewMySQLXMLDump([]byte(sampleXML))
	sql := NewSql(xml, "")
	sql.InsertStmt()
	// Output:
	// INSERT INTO users (id, name, nick_name, state, memo) VALUES (1, "Alice", NULL, 0, "");
	// INSERT INTO orders (id, user_id, item_id, order_date) VALUES (1, 1, 1, "2021-08-09 12:34:45"), (2, 1, 2, "2021-08-09 13:34:45");
}

func ExampleSql_Yaml() {
	xml, _ := NewMySQLXMLDump([]byte(sampleXML))
	sql := NewSql(xml, "")
	sql.Yaml()
	// Output:
	// # users
	// - id: 1
	//   name: "Alice"
	//   nick_name: NULL
	//   state: 0
	//   memo: ""
	//
	// # orders
	// - id: 1
	//   user_id: 1
	//   item_id: 1
	//   order_date: "2021-08-09 12:34:45"
	//
	// - id: 2
	//   user_id: 1
	//   item_id: 2
	//   order_date: "2021-08-09 13:34:45"
}

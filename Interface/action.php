<?php 

header("Content-Type: text/plain");

$action = (isset($_GET["Action"])) ? $_GET["Action"] : NULL;
$id = (isset($_GET["ID"])) ? $_GET["ID"] : NULL;

$db = new SQLite3('../mysqlite_3');

if ($action == "1") {
	$db->query("DELETE FROM IP_table WHERE id = '$id'");
	echo "ok";
} elseif ($action == "2") {
	$db->query("UPDATE IP_table SET release_date=datetime(CURRENT_TIMESTAMP) WHERE id = '$id'");
	echo "ok";
} elseif ($action == "3") {
	$db->query("DELETE FROM IP_table");
	echo "ok";
} elseif ($action == "4") {
	$db->query("INSERT INTO IP_table (id, AddressIP, MAC) VALUES (NULL, '$id', '')");
	echo "ok";
}

?>
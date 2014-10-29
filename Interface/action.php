<?php 

header("Content-Type: text/plain");

function cidrToRange($cidr) {
	$range = array();
	$cidr = explode('/', $cidr);
	$range[0] = long2ip((ip2long($cidr[0])) & ((-1 << (32 - (int)$cidr[1]))));
	for ($i = 1; $i <= (pow(2, (32 - (int)$cidr[1])) - 1); $i++) {
		$range[$i] = long2ip((ip2long($cidr[0])) + $i);
	}
	return $range;
}

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
	$range = cidrToRange(id);
	echo $range[0];
	// $db->query("INSERT INTO IP_table (id, AddressIP, MAC) VALUES (NULL, '$id', '')");
	// echo "ok";
}

?>
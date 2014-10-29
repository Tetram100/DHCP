<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<title>DHCP Server - Admin Interface</title>
	<link rel="icon" type="image/x-icon" href="Assets/favicon.ico" />
	<script src="Assets/jquery-1.11.1.min.js"></script>
	<link href="Assets/bootstrap.min.css" rel="stylesheet">
	<script src="Assets/bootstrap.min.js"></script>
	<script type="text/javascript" src="Assets/oXHR.js"></script>
	<script type="text/javascript">
		<!-- 

		function request(id, action) {
			var xhr = getXMLHttpRequest();

			xhr.onreadystatechange = function() {
				if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 0)) {
					alert(xhr.responseText);
					location.reload();
				}
			};


			xhr.open("GET", "action.php?Action=" + action + "&ID=" + id, true);
			xhr.send(null);
			
		}

		function newIP() {
			var network = encodeURIComponent(document.getElementById("new_network").value);
			var xhr = getXMLHttpRequest();

			xhr.onreadystatechange = function() {
				if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 0)) {
					alert(xhr.responseText);
					location.reload();
				}
			};


			xhr.open("GET", "action.php?Action=4" + "&ID=" + network, true);
			xhr.send(null);

		}

		//-->
	</script>
</head>
<body>
	<div class = "container well">
		<h1>DHCP Project - Admin Page</h1>
		<hr>
		<input type="button" class="btn btn-primary" onclick="request('.$row[0].', 3);" value="Remove all the IP addresses"/>
		<input class="btn btn-primary btn-success" data-toggle="modal" data-target="#newIP" value="Add IP addresses" />
		<hr>
		<table class="table table-striped" style="font-size: 16px;">
			<thead>
				<tr>
					<th>IP addresses</th>
					<th>MAC</th>
					<th>Release Date</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				<?php 
				$db = new SQLite3('../mysqlite_3');
				$results = $db->query("SELECT id, AddressIP, MAC, release_date FROM IP_table");
				while ($row = $results->fetchArray()) {
					echo "<tr>";
					echo "<td>".$row[1]."</td>";
					echo "<td>".substr($row[2], 0, 17)."</td>";
					echo "<td>".$row[3]."</td>";
					echo "<td>";
					echo '<div class="btn-group">';
					echo '<input type="button" class="btn btn-primary" onclick="request('.$row[0].', 1);" value="Remove"/>';
					echo '<input type="button" class="btn btn-success" onclick="request('.$row[0].', 2);" value="Release"/>';
					echo "</div>";
					echo "</td>";
					echo "</tr>";

				}
				?>
			</tbody>
		</table>
	</div>

	<div class="modal fade" id="newIP" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
					<h4 class="modal-title" id="myModalLabel">Add a new network</h4>
				</div>
				<div class="modal-body">
					<p>Write your network with the CIDR notation</p>
					<input type="text" id="new_network" class="form-control" placeholder="Text input">
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
					<input type="button" class="btn btn-primary" onclick="newIP();" value="Add this network"/>
				</div>
			</div>
		</div>
	</div>
</body>
</html>
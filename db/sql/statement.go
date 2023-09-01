// Package sql defines a series of SQL queries and SQL-related constants for interacting with databases,
// particularly MySQL and PostgreSQL. The code is organized into sections for MySQL and PostgreSQL,
// with each section containing relevant SQL queries and related constants.
package sql

const (
	/*------------------------
	 === Common Constants ===
	--------------------------*/
	SQLSelectAll string = `SELECT * FROM %s.%s`
	SQLUpdateRow string = `UPDATE %s SET %s = %s WHERE %s = %s`

	/*------------------------
	 === MySQL Constants ===
	--------------------------*/
	MySQLShowCreateTable   string = `SHOW CREATE TABLE %s.%s`
	MySQLGetColumnDataType        = `
		SELECT 
		    DATA_TYPE
		FROM 
		    INFORMATION_SCHEMA.COLUMNS
		WHERE 
		    TABLE_SCHEMA = '%s'
		AND 
		    TABLE_NAME = '%s'
		AND 
		    COLUMN_NAME = '%s';
	`
	MySQLSchemaSize string = `
		SELECT table_schema "database", 
			sum(data_length + index_length)/1024/1024 "size in MB" 
		FROM 
			information_schema.TABLES 
		WHERE table_schema = '%s' GROUP BY table_schema;
	`
	MySQLShowDatabases     string = `SHOW DATABASES`
	MySQLCountTableColumns string = `
		SELECT 
			count(*) Total_Coulmns 
		FROM 
			information_schema.columns 
		WHERE 
			table_schema = '%s' 
		AND 
			table_name = '%s';
	`
	MySQLCountTableRows string = `SELECT COUNT(*) FROM %s.%s`
	MySQLShowTables     string = `SHOW TABLES`
	MySQLDropTable      string = `DROP TABLE %s`
	MySQLDropDatabase   string = `DROP DATABASE %s`
	MySQLCreateDatabase string = `CREATE DATABASE %s`
	MySQLTruncateTable  string = `TRUNCATE TABLE %s`
	MySQLUse            string = `USE %s`
	MySQLColumnsInfo    string = `
		SELECT
    		c.COLUMN_NAME AS 'Field',
    		c.COLUMN_TYPE AS 'Type',
    		c.COLUMN_KEY AS 'Key',
    		COALESCE(k.CONSTRAINT_NAME, '') AS 'ConstraintName',
    		COALESCE(k.REFERENCED_TABLE_NAME, '') AS 'ReferencedTable',
    		COALESCE(k.REFERENCED_COLUMN_NAME, '') AS 'ReferencedColumn'
		FROM
    		INFORMATION_SCHEMA.COLUMNS c
    	LEFT JOIN 
    		INFORMATION_SCHEMA.KEY_COLUMN_USAGE k 
		ON 
			c.TABLE_NAME = k.TABLE_NAME 
		AND 
			c.COLUMN_NAME = k.COLUMN_NAME
		WHERE
    		c.TABLE_SCHEMA = '%s'
    	AND 
			c.TABLE_NAME = '%s'
	`
	MySQLSelectAllWithLimit string = `SELECT %s FROM %s.%s LIMIT %d OFFSET %d`
	MySQLGetTablesSize      string = `
		SELECT
			TABLE_NAME AS "Table",
			ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS "Size (MB)"
		FROM
			information_schema.TABLES
		WHERE
			TABLE_SCHEMA = '%s'
		ORDER BY
			(DATA_LENGTH + INDEX_LENGTH) DESC;
	`
	MySQLGetTableSize string = `
		SELECT
			table_name AS "Table",
			ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS "Size (MB)"
		FROM
			information_schema.TABLES
		WHERE
			TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s';
	`

	/*---------------------------
	 === PostgreSQL Constants ===
	-----------------------------*/
	PostgreSQLShowDatabases string = `
		SELECT
			datname
		FROM
			pg_database
		WHERE NOT 
			datistemplate
	`
	PostgreSQLShowTables string = `
		SELECT 
			table_name 
		FROM 
			information_schema.tables 
		WHERE 
			table_schema = '%s'
	`
	PostgreSQLSelectAllWithLimit string = `SELECT %s FROM %s.%s LIMIT %d OFFSET %d`
	PostgreSQLSchemaSize         string = `
		SELECT 
			pg_size_pretty(pg_database_size(current_database())) 
		AS 
			"database size"
	`
	PostgreSQLCountTableColumns string = `
		SELECT 
			count(column_name) AS "Total_Columns"
		FROM 
			information_schema.columns
		WHERE 
			table_schema = '%s'
  		AND 
			table_name = '%s'
	`
	PostgreSQLCountTableRows string = `
		SELECT 
			count(*) AS "Total_Rows" 
		FROM 
			%s.%s
	`
	PostgreSQLTableSize string = `
		WITH table_info AS (
    		SELECT
        		'%s' AS schema_name,
        		'%s' AS table_name
		)
		SELECT
    		table_info.table_name 
		AS "Table_Name",    
			ROUND(((pg_total_relation_size(table_info.schema_name || '.' || table_info.table_name)) / 1024.0 / 1024.0), 2)
		AS "Table_Size"
		FROM
			table_info;
	`
	PostgreSQLGetColumnDataType = `
		SELECT 
		    data_type
		FROM 
		    information_schema.columns
		WHERE 
		    table_schema = '%s' 
		AND 
			table_name = '%s' 
		AND
			column_name = '%s';
	`
	PostgreSQLTableSizes string = `
		SELECT 
			table_name 
		AS "Table",
       		pg_size_pretty(pg_total_relation_size('"' || table_schema || '"."' || table_name || '"')) 
		AS "Table_Size"
		FROM 
			information_schema.tables
		WHERE 
			table_type = 'BASE TABLE'
      	AND 
			table_schema 
		NOT IN 
			('pg_catalog', 'information_schema')
	`
	PostgreSQLDropTable      string = `DROP TABLE IF EXISTS %s`
	PostgreSQLDropDatabase   string = `DROP DATABASE IF EXISTS %s`
	PostgreSQLCreateDatabase string = `CREATE DATABASE %s`
	PostgreSQLTruncateTable  string = `TRUNCATE TABLE %s`
	PostgreSQLColumnsInfo    string = `
		SELECT 
			c.column_name AS Field, 
			c.data_type AS Type,
			CASE
				WHEN tc.constraint_type = 'PRIMARY KEY' THEN 'PRI'
				WHEN tc.constraint_type = 'FOREIGN KEY' THEN 'MUL'
				ELSE '' 
				END AS Key,
				COALESCE(tc.constraint_name, '') AS ConstraintName,
				COALESCE(ccu.table_name, '') AS ReferencedTable,
				COALESCE(ccu.column_name, '') AS ReferencedColumn
			FROM 
				information_schema.columns c
			LEFT JOIN 
				information_schema.key_column_usage kcu 
			ON 
				c.table_name = kcu.table_name AND c.column_name = kcu.column_name
			LEFT JOIN 
				information_schema.table_constraints tc 
			ON 
				kcu.constraint_name = tc.constraint_name
			LEFT JOIN 
				information_schema.constraint_column_usage ccu 
			ON 
				tc.constraint_name = ccu.constraint_name
			WHERE 
				c.table_schema = '%s' 
			AND 
				c.table_name = '%s';
	`

	// PostgreSQLShowCreateFunction is function that attempts to resemble the behaviour of mysql's 'show create' statement
	// the code is taken from an old answer found in
	// https://stackoverflow.com/questions/2593803/how-to-generate-the-create-table-sql-statement-for-an-existing-table-in-postgr
	PostgreSQLShowCreateFunction = `
		CREATE OR REPLACE FUNCTION public.show_create_table(
			in_schema_name varchar,
			in_table_name varchar
		)
		RETURNS text
		LANGUAGE plpgsql VOLATILE
		AS
		$$
  		DECLARE
    		-- the ddl we're building
    		v_table_ddl text;
			-- data about the target table
			v_table_oid int;
			-- records for looping
			v_column_record record;
			v_constraint_record record;
			v_index_record record;
  		BEGIN
			SELECT 
				c.oid INTO v_table_oid
			FROM 
				pg_catalog.pg_class c
			LEFT JOIN 
				pg_catalog.pg_namespace n 
			ON 
				n.oid = c.relnamespace
    	WHERE 1=1
      		AND 
				c.relkind = 'r'
			AND 
				c.relname = in_table_name -- the table name
			AND 
				n.nspname = in_schema_name; -- the schema

		-- start the create definition
		v_table_ddl := 'CREATE TABLE ' || in_schema_name || '.' || in_table_name || ' (' || E'\n';

		FOR v_column_record IN
		SELECT
			c.column_name,
			c.data_type,
			c.character_maximum_length,
			c.is_nullable,
			c.column_default
		FROM 
			information_schema.columns c
		WHERE 
			(table_schema, table_name) = (in_schema_name, in_table_name)
		ORDER BY 
			ordinal_position
		LOOP
		v_table_ddl := v_table_ddl || '  ' -- note: two char spacer to start, to indent the column
			|| v_column_record.column_name || ' '
			|| v_column_record.data_type || CASE WHEN v_column_record.character_maximum_length IS NOT NULL THEN ('(' || v_column_record.character_maximum_length || ')') ELSE '' END || ' '
			|| CASE WHEN v_column_record.is_nullable = 'NO' THEN 'NOT NULL' ELSE 'NULL' END
			|| CASE WHEN v_column_record.column_default IS NOT null THEN (' DEFAULT ' || v_column_record.column_default) ELSE '' END
			|| ',' || E'\n';
		END LOOP;

		FOR v_constraint_record IN
		SELECT
			con.conname as constraint_name,
			con.contype as constraint_type,
			CASE
			WHEN 
				con.contype = 'p' THEN 1 -- primary key constraint
			WHEN 
				con.contype = 'u' THEN 2 -- unique constraint
			WHEN 
				con.contype = 'f' THEN 3 -- foreign key constraint
			WHEN 
				con.contype = 'c' THEN 4
			ELSE 5
			END as type_rank,
			pg_get_constraintdef(con.oid) as constraint_definition
		FROM 
			pg_catalog.pg_constraint con
		JOIN 
			pg_catalog.pg_class rel ON rel.oid = con.conrelid
		JOIN 
			pg_catalog.pg_namespace nsp ON nsp.oid = connamespace
		WHERE 
			nsp.nspname = in_schema_name
		AND 
			rel.relname = in_table_name
		ORDER BY type_rank
		LOOP
		v_table_ddl := v_table_ddl || '  ' -- note: two char spacer to start, to indent the column
			|| 'CONSTRAINT' || ' '
			|| v_constraint_record.constraint_name || ' '
			|| v_constraint_record.constraint_definition
			|| ',' || E'\n';
		END LOOP;

		-- drop the last comma before ending the create statement
		v_table_ddl = substr(v_table_ddl, 0, length(v_table_ddl) - 1) || E'\n';

		-- end the create definition
		v_table_ddl := v_table_ddl || ');' || E'\n';

		-- suffix create statement with all of the indexes on the table
		FOR v_index_record IN
			SELECT indexdef
			FROM pg_indexes
			WHERE (schemaname, tablename) = (in_schema_name, in_table_name)
		LOOP
			v_table_ddl := v_table_ddl
				|| v_index_record.indexdef
				|| ';' || E'\n';
		END LOOP;

    	-- return the ddl
    	RETURN v_table_ddl;
  		END;
		$$;
	`
	PostgreSQLShowCreate             = `SELECT * FROM public.show_create_table('%s', '%s');`
	PostgreSQLDropShowCreateFunction = `DROP FUNCTION public.show_create_table(varchar, varchar);`
)

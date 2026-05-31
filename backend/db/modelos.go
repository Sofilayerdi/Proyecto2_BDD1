package db

type Producto struct {
	IDProducto  int     `gorm:"primaryKey;column:id_producto"`
	Nombre      string  `gorm:"nombre"`
	Categoria   string  `gorm:"categoria"`
	IDProveedor int     `gorm:"id_proveedor"`
	Cantidad    int     `gorm:"cantidad"`
	Precio      float64 `gorm:"precio"`
}

func (Producto) TableName() string { return "producto" }

type Ramo struct {
	IDRamo int     `gorm:"primaryKey;column:id_ramo"`
	Total  float64 `gorm:"total"`
}

func (Ramo) TableName() string { return "ramo" }

type RamoProducto struct {
	IDRamo     int `gorm:"primaryKey;column:id_ramo"`
	IDProducto int `gorm:"primaryKey;column:id_producto"`
	Cantidad   int `gorm:"column:cantidad"`
}

func (RamoProducto) TableName() string { return "ramo_producto" }

type Venta struct {
	IDVenta     int     `gorm:"primaryKey;column:id_venta"`
	IDCliente   int     `gorm:"column:id_cliente"`
	IDEmpleado  int     `gorm:"column:id_empleado"`
	Fecha       string  `gorm:"fecha"`
	PrecioTotal float64 `gorm:"precio_total"`
}

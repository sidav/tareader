package cob

// Values and comments taken from EXPTYPE.H in game scrips

// Constants unsed in the COB machine SetValue and GetValue opcodes.
const (
	ValueId_ACTIVATION         = 1  // set or get
	ValueId_STANDINGMOVEORDERS = 2  // set or get
	ValueId_STANDINGFIREORDERS = 3  // set or get
	ValueId_HEALTH             = 4  // get (0-100%)
	ValueId_INBUILDSTANCE      = 5  // set or get
	ValueId_BUSY               = 6  // set or get (used by misc. special case missions like transport ships)
	ValueId_PIECE_XZ           = 7  // get
	ValueId_PIECE_Y            = 8  // get
	ValueId_UNIT_XZ            = 9  // get
	ValueId_UNIT_Y             = 10 // get
	ValueId_UNIT_HEIGHT        = 11 // get
	ValueId_XZ_ATAN            = 12 // get atan of packed x,z coords
	ValueId_XZ_HYPOT           = 13 // get hypot of packed x,z coords
	ValueId_ATAN               = 14 // get ordinary two-parameter atan
	ValueId_HYPOT              = 15 // get ordinary two-parameter hypot
	ValueId_GROUND_HEIGHT      = 16 // get
	ValueId_BUILD_PERCENT_LEFT = 17 // get 0 = unit is built and ready, 1-100 = How much is left to build
	ValueId_YARD_OPEN          = 18 // set or get (change which plots we occupy when building opens and closes)
	ValueId_BUGGER_OFF         = 19 // set or get (ask other units to clear the area)
	ValueId_ARMORED            = 20 // set or get
)

// SFX types?
const (
	EffectId_SHATTER        = 1  // The piece will shatter instead of remaining whole
	EffectId_EXPLODE_ON_HIT = 2  // The piece will explode when it hits the ground
	EffectId_FALL           = 4  // The piece will fall due to gravity instead of just flying off
	EffectId_SMOKE          = 8  // A smoke trail will follow the piece through the air
	EffectId_FIRE           = 16 // A fire trail will follow the piece through the air
	EffectId_BITMAPONLY     = 32 // The piece will not fly off or shatter or anything.  Only a bitmap explosion will be rendered.
)

// Bitmap explosion types
const (
	BitmapExpType_BITMAP1    = 256
	BitmapExpType_BITMAP2    = 512
	BitmapExpType_BITMAP3    = 1024
	BitmapExpType_BITMAP4    = 2048
	BitmapExpType_BITMAP5    = 4096
	BitmapExpType_BITMAPNUKE = 8192

	BitmapExpType_BITMAPMASK = 16128 // Mask of the possible bitmap bits
)

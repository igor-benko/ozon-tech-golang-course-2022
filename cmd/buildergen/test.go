package main

type (
	testInnerStruct struct{}
	TestInnerStruct struct{}
	testStruct      struct{}
	TestStruct      struct {
		PublicValueInt           int
		PublicValueString        string
		PublicValueFloat32       float32
		PublicValuePublicStruct  TestInnerStruct
		PublicValuePrivateStruct testInnerStruct

		PublicPointerInt           *int
		PublicPointerString        *string
		PublicPointerFloat32       *float32
		PublicPointerPublicStruct  *TestInnerStruct
		PublicPointerPrivateStruct *testInnerStruct

		privateValueInt           int
		privateValueString        string
		privateValueFloat32       float32
		privateValuePublicStruct  TestInnerStruct
		privateValuePrivateStruct testInnerStruct

		privatePointerInt           *int
		privatePointerString        *string
		privatePointerFloat32       *float32
		privatePointerPublicStruct  *TestInnerStruct
		privatePointerPrivateStruct *testInnerStruct
	}
)

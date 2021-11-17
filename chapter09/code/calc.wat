(module
  (type $ft0 (func (param i64 i64) (result i64)))
  (table funcref (elem $add $sub $mul))

  (func $calc (param $a i64) (param $b i64) (param $op i64) (result i64)
    (local.get $a) (local.get $b) (local.get $op)
    (call_indirect (type $ft0))
  )

  (func $add (type $ft0) (i64.add (local.get 0) (local.get 1)))
  (func $sub (type $ft0) (i64.sub (local.get 0) (local.get 1)))
  (func $mul (type $ft0) (i64.mul (local.get 0) (local.get 1)))
)
# Handle DB Transaction

This guide will help you fix the issue regard to rollback/commit on DB transaction, making some changes on services layer will help your db transaction work properly.

### There are 2 solutions to do transaction

## Solution 1: Do transaction manually

from `app/service/sample_product_service.go`

1. You have to explicitly begin a new transaction (`tx`)
```golang
// s.baseRepo.BeginTx()
tx := s.baseRepo.GetDB().Begin()
```

2. Use `tx` to create/update/delete instead of `s.sampleProductRepo`
```golang
// result, err := s.sampleProductRepo.Create(&product)
err := tx.Create(&product).Error
```

3. Use `tx` to rollback instead of `s.baseRepo`
```golang
// s.baseRepo.RollbackTx()
tx.Rollback()
```

4. Use `tx` to commit instead of `s.baseRepo`
```golang
// s.baseRepo.CommitTx()
tx.Commit()
```

### Final code snippet
```golang
logger := log.With(s.logger, "SampleProductService", "CreateSampleProduct")
tx := s.baseRepo.GetDB().Begin()
//Set request to entity
product := entity.SampleProduct{
    Name:   input.Name,
    Sku:    input.Sku,
    Uom:    input.Uom,
    Weight: input.Weight,
    // Set Jwt Info to Entity
    BaseIDModel: base.BaseIDModel {
        CreatedByUid: input.ActorUID,
        CreatedBy:    input.ActorName,
    },
}

err := tx.Create(&product).Error
if err != nil {
    _ = level.Error(logger).Log(err)
    tx.Rollback()
    return nil, message.FailedMsg
}
tx.Commit()

return response.SampleProductMapToResponse(product), message.SuccessMsg
```
This approach will help your DB transaction work properly.

## Solution 2: Short way

This solution will do beginTx/commitTx/rollbackTx implicitly

1. Get `db` from Repository
```golang
// s.baseRepo.BeginTx()
db := s.baseRepo.GetDB()
```
2. Perform a set of operations within a `Transaction`
```golang
// result, err := s.sampleProductRepo.Create(&product)
err := db.Transaction(func(tx *gorm.DB) error {
    return tx.Create(&product).Error
})
```
3. Remove explicit commit/rollback
```golang
// s.baseRepo.RollbackTx()
// s.baseRepo.CommitTx()
```

### Final code snippet
```golang
logger := log.With(s.logger, "SampleProductService", "CreateSampleProduct")
db := s.baseRepo.GetDB()
//Set request to entity
product := entity.SampleProduct{
    Name:   input.Name,
    Sku:    input.Sku,
    Uom:    input.Uom,
    Weight: input.Weight,
    // Set Jwt Info to Entity
    BaseIDModel: base.BaseIDModel{
        CreatedByUid: input.ActorUID,
        CreatedBy:    input.ActorName,
    },
}

err := db.Transaction(func(tx *gorm.DB) error {
    return tx.Create(&product).Error
})

if err != nil {
    _ = level.Error(logger).Log(err)
    return nil, message.FailedMsg
}

return response.SampleProductMapToResponse(product), message.SuccessMsg
```

## The better way (new project only)
For new project you could fork from branch `refactor/repository-layer` which already fixed the db transaction issue
- ### Short way
You could use transaction with simple way on repository layer (e.g: `app/repository/sample_product_repository.go`)
```golang
err := r.BaseRepository.Transaction(func(tx *gorm.DB) error {
    return tx.Create(product).Error
})
```

- ### Manual way
This solution allow you to controll transaction manually
```golang
tx := r.BaseRepository.BeginTx()

if err := tx.Create(product).Error; err != nil {
	tx.Rollback()
}

tx.Commit()
```

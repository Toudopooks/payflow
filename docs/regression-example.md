# Regression Example: Duplicate Payment Retry

## Problem

A client can lose the response after a payment commits and retry the same request. A naive handler creates a second payment attempt and may charge the installment twice.

## Root Cause

Checking installment status before starting the write transaction leaves a race. Two requests can both observe `pending`, then both proceed.

## Regression Test

The MySQL integration test sends the same payment command twice with one idempotency key. It asserts that the second call is marked as a replay and returns the first payment ID. It then sends a new key for the paid installment and expects `ErrInstallmentAlreadyPaid`.

## Fix

The repository inserts or locks the idempotency record and locks the installment with `SELECT ... FOR UPDATE`. Payment creation, installment state, event insertion, and idempotency completion occur in the same transaction.

## Result

Identical retries return the original resource. A key reused for different input returns a conflict. A new request cannot pay an already completed installment.


package main

import (
	"errors"
	"fmt"
	"time"
)

//Account is an interface that wraps the common behavior for accounts.
type Account interface {
	Deposit(amount Money) error
	Withdraw(amount Money) error
}

//Transfer transfers money from one account to another.
//Returns the root cause in case of an error.
func Transfer(fromAcct Account, toAcct Account, amount Money) error {
	if err := fromAcct.Withdraw(amount); err == nil {
		if depErr := toAcct.Deposit(amount); depErr == nil {
			fmt.Printf("Transfered %f from %s to %+v", amount, fromAcct, toAcct)
		} else {
			//return the root cause
			return depErr
		}
	} else {
		//return the root cause
		return err
	}
	return nil
}

//Money type.
type Money float64

//SavingsAccount encapsulates the state of a savings account.
type SavingsAccount struct {
	InterestRate float32
	MinBalance   Money

	Num      string
	Name     string
	OpenDate time.Time
	Balance  Money
}

//CheckingAccount encapsulates the state of a checking account.
type CheckingAccount struct {
	TransactionFee   float32
	OverDraftEnabled bool

	Num      string
	Name     string
	OpenDate time.Time
	Balance  Money
}

//OpenSavingsAccount is the Initializing constructor for SavingsAccount
func OpenSavingsAccount(no string, name string, openingDate time.Time) *SavingsAccount {
	a := SavingsAccount{
		Num:          no,
		Name:         name,
		OpenDate:     openingDate,
		InterestRate: 0.9,
		MinBalance:   15.0,
	}
	return &a
}

//OpenCheckingAccount is the initializing constructor for CheckingAccount
func OpenCheckingAccount(no string, name string, openingDate time.Time, overdraftFlag bool) *CheckingAccount {
	return &CheckingAccount{
		Num:              no,
		Name:             name,
		OpenDate:         openingDate,
		TransactionFee:   0.15,
		OverDraftEnabled: overdraftFlag,
	}
}

//Deposit adds specified amount to existing balance of savings account.
//Returns nil on successful deposit.
func (acct *SavingsAccount) Deposit(amount Money) error {
	fmt.Printf("Depositing %f \n", amount)
	acct.Balance = acct.Balance + amount
	return nil
}

//Withdraw removes specified amount from existing balance of savings account.
//Returns nil on successful withdrawal.
func (acct *SavingsAccount) Withdraw(amount Money) error {
	//check for min balance invariant
	if acct.Balance-amount < acct.MinBalance {
		return errors.New("Not enough money to withdraw in SavingsAccount")
	}
	acct.Balance = acct.Balance - amount
	return nil
}

//Deposit adds specified amount to existing balance of checking account.
//Returns nil on successful deposit.
func (acct *CheckingAccount) Deposit(amount Money) error {
	fmt.Printf("Depositing %f \n", amount)
	acct.Balance = acct.Balance + amount
	return nil
}

//Withdraw removes specified amount from existing balance of checking account.
//Returns nil on successful withdrawal.
func (acct *CheckingAccount) Withdraw(amount Money) error {
	//check the overdraft invariant
	if acct.Balance < amount && !acct.OverDraftEnabled {
		return errors.New("Not enough money to withdraw in CheckingAccount")
	}
	acct.Balance = acct.Balance - amount
	return nil
}

func main() {
	aliceAcct := OpenSavingsAccount("12345", "Alice", time.Date(1999, time.January, 03, 0, 0, 0, 0, time.UTC))
	fmt.Println("Alice's account =", aliceAcct)
	aliceAcct.Deposit(Money(100.0))
	fmt.Println("Alice's account (after deposit) =", aliceAcct)
	if err := aliceAcct.Withdraw(Money(10)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Alice's account (after withdrawl) =", aliceAcct)
	}

	const overDraft = false
	bobAcct := OpenCheckingAccount("98765", "Bob", time.Date(1997, time.April, 03, 0, 0, 0, 0, time.UTC), overDraft)
	fmt.Println("\nBob's account =", bobAcct)
	bobAcct.Deposit(Money(100.0))
	fmt.Println("Bob's account (after deposit) =", bobAcct)
	if err := bobAcct.Withdraw(Money(77)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Bob's account (after withdrawl) =", bobAcct)
	}

	fmt.Println("\nTransfering 76 from Alice to Bob's acct")
	if err := Transfer(aliceAcct, bobAcct, Money(76.0)); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Alice's account =", aliceAcct)
	fmt.Println("Bob's account =", bobAcct)
}

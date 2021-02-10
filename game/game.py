import sys
from random import randrange


class Game:

    @staticmethod
    def start():
        number = randrange(1000)
        tries = 0

        while True:
            print("Please enter a guess\n")
            guess = input("Number: ")

            if guess.isdigit():
                if Game.guess(int(guess), number):
                    break
            else:
                print("Enter a real number please!")

            tries += 1

        print("Well done, you have found the number. Amount of tries: " + str(tries))

    @staticmethod
    def guess(guess, number):
        if guess == number:
            return True
        elif guess > number:
            print("Guess was higher than number")
            return False
        else:
            print("Guess was lower than number")
            return False


if __name__ == '__main__':
    Game.start()

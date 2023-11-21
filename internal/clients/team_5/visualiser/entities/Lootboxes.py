"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import pygame
import pygame_gui
from visualiser.util.Constants import LOOTBOX, OVERLAY
from visualiser.entities.Common import Drawable

class Lootbox(Drawable):
    def __init__(self, x:int, y:int, lootboxID) -> None:
        super().__init__(x, y)
        self.id = lootboxID
        self.colour = LOOTBOX["DEFAULT_COLOUR"]

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the lootbox
        """
        # Determine the grid size
        self.trueX = int(self.x * zoom + offsetX)
        self.trueY = int(self.y * zoom + offsetY)
        # Draw the lootbox
        border = pygame.Surface(((2*LOOTBOX["LINE_WIDTH"] + LOOTBOX["WIDTH"])*zoom, (2*LOOTBOX["LINE_WIDTH"] + LOOTBOX["HEIGHT"])*zoom))
        border.fill(LOOTBOX["LINE_COLOUR"])
        overlay = pygame.Surface((LOOTBOX["WIDTH"]*zoom, LOOTBOX["HEIGHT"]*zoom))
        overlay.fill(self.colour)
        # Add lootbox text
        font = pygame.font.SysFont(OVERLAY["FONT"], int(LOOTBOX["FONT_SIZE"] * zoom))
        if self.colour in ("White", "Yellow"):
            text = font.render(self.id, True, "Black")
        else:
            text = font.render(self.id, True, "White")
        # center the text
        textX = (LOOTBOX["WIDTH"]*zoom - text.get_width()) / 2
        textY = (LOOTBOX["HEIGHT"]*zoom - text.get_height()) / 2
        overlay.blit(text, (textX, textY))
        # add the overlay to the border
        border.blit(overlay, (LOOTBOX["LINE_WIDTH"]*zoom, LOOTBOX["LINE_WIDTH"]*zoom))
        screen.blit(border, (self.trueX, self.trueY))
        # Draw the agents within the bike
        self.overlay = self.update_overlay(zoom)
        self.draw_overlay(screen)

    def check_collision(self, mouseX: int, mouseY: int, offsetX: int, offsetY: int, zoom: float) -> bool:
        """
        Check if the mouse click intersects with the bike.
        """
        return (self.trueX <= mouseX <= self.trueX + LOOTBOX["WIDTH"]) and \
               (self.trueY <= mouseY <= self.trueY + LOOTBOX["HEIGHT"])

    def change_round(self, json:dict) -> None:
        """
        Change the current round for the agents
        """
        self.colour = json[self.id]["colour"]
        self.properties = {
            "Position" : f"{json[self.id]['position']['x']}, {json[self.id]['position']['y']}",
        }

    def propagate_click(self, mouseX:int, mouseY:int, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Propagate the click
        """
        self.click(mouseX, mouseY, offsetX, offsetY, zoom)

using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace maktba.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class InitialCreate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "Authors",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    Name = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: false),
                    Bio = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Authors", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Sources",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    Title = table.Column<string>(type: "character varying(500)", maxLength: 500, nullable: false),
                    Type = table.Column<int>(type: "integer", nullable: false),
                    Description = table.Column<string>(type: "text", nullable: true),
                    CoverUrl = table.Column<string>(type: "text", nullable: true),
                    Url = table.Column<string>(type: "text", nullable: true),
                    PublishedDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Sources", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Taxonomies",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    Name = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    Type = table.Column<int>(type: "integer", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Taxonomies", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "SourceAuthors",
                columns: table => new
                {
                    AuthorsId = table.Column<Guid>(type: "uuid", nullable: false),
                    SourcesId = table.Column<Guid>(type: "uuid", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_SourceAuthors", x => new { x.AuthorsId, x.SourcesId });
                    table.ForeignKey(
                        name: "FK_SourceAuthors_Authors_AuthorsId",
                        column: x => x.AuthorsId,
                        principalTable: "Authors",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_SourceAuthors_Sources_SourcesId",
                        column: x => x.SourcesId,
                        principalTable: "Sources",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "SourceTaxonomies",
                columns: table => new
                {
                    SourcesId = table.Column<Guid>(type: "uuid", nullable: false),
                    TaxonomiesId = table.Column<Guid>(type: "uuid", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_SourceTaxonomies", x => new { x.SourcesId, x.TaxonomiesId });
                    table.ForeignKey(
                        name: "FK_SourceTaxonomies_Sources_SourcesId",
                        column: x => x.SourcesId,
                        principalTable: "Sources",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_SourceTaxonomies_Taxonomies_TaxonomiesId",
                        column: x => x.TaxonomiesId,
                        principalTable: "Taxonomies",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_SourceAuthors_SourcesId",
                table: "SourceAuthors",
                column: "SourcesId");

            migrationBuilder.CreateIndex(
                name: "IX_Sources_Title",
                table: "Sources",
                column: "Title");

            migrationBuilder.CreateIndex(
                name: "IX_SourceTaxonomies_TaxonomiesId",
                table: "SourceTaxonomies",
                column: "TaxonomiesId");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "SourceAuthors");

            migrationBuilder.DropTable(
                name: "SourceTaxonomies");

            migrationBuilder.DropTable(
                name: "Authors");

            migrationBuilder.DropTable(
                name: "Sources");

            migrationBuilder.DropTable(
                name: "Taxonomies");
        }
    }
}
